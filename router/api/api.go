package api

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
)

func RegisterRoutes(m *macaron.Macaron) {
	m.Group("/auth", func() {
		m.Post("/", binding.Json(AuthRequest{}), Auth)
		//m.Post("/", api.AuthCookies)
	})
}
