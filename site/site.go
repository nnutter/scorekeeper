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
			"<style>" + ui.CSS() + "</style>",
		},
	}
}
