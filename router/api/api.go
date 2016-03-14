package api

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
)

type Response struct {
	Code int
}

type ErrorResponse struct {
	Response
	Error string
	ErrorDescription string
}

func RegisterRoutes(m *macaron.Macaron) {
	m.Group("/auth", func() {
		m.Post("/", binding.Json(AuthRequest{}), Auth)
		m.Post("/cookies", binding.Json(AuthCookiesRequest{}), AuthCookies)
	})

	m.Post("/bugs/crash", binding.Json(CrashRequest{}), Crash)
}
