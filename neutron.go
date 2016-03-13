package main

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/router/api"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	m.Group("/api", func() {
		api.RegisterRoutes(m)
	})

	m.Use(macaron.Static("public/build", macaron.StaticOptions{
		IndexFile: "app.html",
		SkipLogging: true,
	}))

	m.Run()
}
