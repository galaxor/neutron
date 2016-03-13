package main

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/router/api"
)

func main() {
	m := macaron.Classic()

	m.Group("/api", func() {
		api.RegisterRoutes(m)
	})

	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Run()
}
