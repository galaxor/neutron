package main

import (
	"io/ioutil"
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend/memory"
	"github.com/emersion/neutron/router/api"
)

func main() {
	publicDir := "public/build"
	indexFile := "app.html"

	backend := memory.New()

	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// API
	m.Group("/api", func() {
		api.New(m, backend)
	})

	// Serve static files
	m.Use(macaron.Static(publicDir, macaron.StaticOptions{
		IndexFile: indexFile,
		SkipLogging: true,
	}))

	m.NotFound(func(ctx *macaron.Context) {
		// API endpoint, send error
		if strings.HasPrefix(ctx.Req.URL.Path, "/api/") {
			ctx.PlainText(404, []byte("endpoint not found"))
			return
		}

		// Fallback to index file

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
