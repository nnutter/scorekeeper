package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/nnutter/scorekeeper/ui"
)

func main() {
	app.Route("/", func() app.Composer {
		return ui.New()
	})
	app.RunWhenOnBrowser()
	if app.IsClient {
		return
	}

	if err := ensureWASM(); err != nil {
		log.Fatal(err)
	}

	h := &app.Handler{
		Name:        "Scorekeeper",
		Description: "Baseball scorekeeping with simplified Retrosheet export",
		Resources:   app.LocalDir("."),
		RawHeaders: []string{
			"<style>" + ui.CSS() + "</style>",
		},
	}

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8000"
	}

	log.Printf("listening on :%s", addr)
	if err := http.ListenAndServe(":"+addr, h); err != nil {
		log.Fatal(err)
	}
}

func ensureWASM() error {
	if err := os.MkdirAll("web", 0o755); err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-o", "web/app.wasm", ".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &buildError{output: string(output), err: err}
	}
	return nil
}

type buildError struct {
	output string
	err    error
}

func (e *buildError) Error() string {
	if e.output == "" {
		return "building web/app.wasm failed: " + e.err.Error()
	}
	return "building web/app.wasm failed: " + e.err.Error() + ": " + e.output
}
