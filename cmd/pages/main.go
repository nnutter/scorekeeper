package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/nnutter/scorekeeper/site"
)

func main() {
	var outputDir string
	var repoName string
	flag.StringVar(&outputDir, "output", "dist", "directory for generated static site")
	flag.StringVar(&repoName, "repo", "", "GitHub repository name for Pages base path")
	flag.Parse()

	if err := os.MkdirAll(filepath.Join(outputDir, "web"), 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, ".nojekyll"), []byte{}, 0o644); err != nil {
		fatal(err)
	}

	site.RegisterRoutes()

	resources := app.ResourceResolver(app.LocalDir(outputDir))
	if repoName != "" {
		resources = app.GitHubPages(repoName)
	}

	if err := app.GenerateStaticWebsite(outputDir, site.NewHandler(resources)); err != nil {
		fatal(err)
	}

	assets := []string{"baseball-icon.svg", "baseball-icon-192.png", "baseball-icon-512.png"}
	iconAssets, err := filepath.Glob(filepath.Join("web", "icon-*.svg"))
	if err != nil {
		fatal(err)
	}
	for _, path := range iconAssets {
		assets = append(assets, filepath.Base(path))
	}

	for _, name := range assets {
		if err := copyFile(filepath.Join("web", name), filepath.Join(outputDir, "web", name)); err != nil {
			fatal(err)
		}
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
