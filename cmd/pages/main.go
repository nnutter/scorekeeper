package main

import (
	"flag"
	"fmt"
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
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
