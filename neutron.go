package main

import (
	"io/ioutil"

	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/router/api"
)

func main() {
	publicDir := "public/build"
	indexFile := "app.html"

	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// API
	m.Group("/api", func() {
		api.RegisterRoutes(m)

		m.NotFound(func(ctx *macaron.Context) {
			ctx.PlainText(404, []byte("endpoint not found"))
		})
	})

	// Serve static files
	m.Use(macaron.Static(publicDir, macaron.StaticOptions{
		IndexFile: indexFile,
		SkipLogging: true,
	}))

	// Fallback to index file
	m.NotFound(func(ctx *macaron.Context) {
		data, err := ioutil.ReadFile(publicDir + "/" + indexFile)
		if err != nil {
			ctx.PlainText(404, []byte("page not found"))
			return
		}

		ctx.Resp.Header().Set("Content-Type", "text/html")
		ctx.Resp.Write(data)
	})

	m.Run()
}
