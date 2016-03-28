package main

import (
	"io/ioutil"

	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/config"
	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/memory"
	"github.com/emersion/neutron/backend/imap"
	"github.com/emersion/neutron/backend/smtp"
	"github.com/emersion/neutron/router/api"
)

const (
	publicDir = "public/build"
	indexFile = "app.html"
)

func main() {
	// Load config
	c, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}

	// Create backend
	bkd := backend.New()
	if c.Memory != nil && c.Memory.Enabled {
		memory.Use(bkd)
		if c.Memory.Populate {
			memory.Populate(bkd)
		}
	}
	if c.Imap != nil && c.Imap.Enabled {
		passwords := imap.Use(bkd, &c.Imap.Config)

		if c.Smtp != nil && c.Smtp.Enabled {
			smtp.Use(bkd, &c.Smtp.Config, passwords)
		}
	}

	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// API
	m.Group("/api", func() {
		api.New(m, bkd)
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
