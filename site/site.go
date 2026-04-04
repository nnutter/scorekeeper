package site

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/nnutter/scorekeeper/ui"
)

func RegisterRoutes() {
	app.Route("/", func() app.Composer {
		return ui.New()
	})
}

func NewHandler(resources app.ResourceResolver) *app.Handler {
	return &app.Handler{
		Name:        "Scorekeeper",
		Description: "Baseball scorekeeping with simplified Retrosheet export",
		Resources:   resources,
		RawHeaders: []string{
			`<link rel="icon" type="image/svg+xml" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 64 64'%3E%3Ccircle cx='32' cy='32' r='28' fill='%23fffdf8' stroke='%23c7b299' stroke-width='2'/%3E%3Cpath d='M22 13c-5 6-8 12-8 19s3 13 8 19' fill='none' stroke='%23b3261e' stroke-width='2.5' stroke-linecap='round'/%3E%3Cpath d='M42 13c5 6 8 12 8 19s-3 13-8 19' fill='none' stroke='%23b3261e' stroke-width='2.5' stroke-linecap='round'/%3E%3Cpath d='M24 18l3 4m-5 2l4 3m-5 3l4 2m-4 5l4-2m-3 6l4-3m-2 8l3-4' fill='none' stroke='%23b3261e' stroke-width='2' stroke-linecap='round'/%3E%3Cpath d='M40 18l-3 4m5 2l-4 3m5 3l-4 2m4 5l-4-2m3 6l-4-3m2 8l-3-4' fill='none' stroke='%23b3261e' stroke-width='2' stroke-linecap='round'/%3E%3C/svg%3E">`,
			"<style>" + ui.CSS() + "</style>",
		},
	}
}
