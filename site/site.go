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
		Icon: app.Icon{
			Default:  "/web/baseball-icon-192.png",
			Large:    "/web/baseball-icon-512.png",
			Maskable: "/web/baseball-icon-512.png",
			SVG:      "/web/baseball-icon.svg",
		},
		RawHeaders: []string{
			"<style>" + ui.CSS() + "</style>",
		},
	}
}
