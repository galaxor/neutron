package api

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"

	"github.com/emersion/neutron/backend"
)

type Api struct {
	backend backend.Backend
}

type Req struct {}

type Resp struct {
	Code int
}

type ErrorResp struct {
	Resp
	Error string
	ErrorDescription string
}

func New(m *macaron.Macaron, backend backend.Backend) {
	api := &Api{
		backend: backend,
	}

	m.Use(func (ctx *macaron.Context) {
		if appVersion, ok := ctx.Req.Header["X-Pm-Appversion"]; ok {
			ctx.Data["appVersion"] = appVersion
		}
		if apiVersion, ok := ctx.Req.Header["X-Pm-Apiversion"]; ok {
			ctx.Data["apiVersion"] = apiVersion
		}
		if sessionToken, ok := ctx.Req.Header["X-Pm-Session"]; ok {
			ctx.Data["sessionToken"] = sessionToken
		}
	})

	m.Group("/auth", func() {
		m.Post("/", binding.Json(AuthReq{}), api.Auth)
		m.Post("/cookies", binding.Json(AuthCookiesReq{}), api.AuthCookies)
	})

	m.Group("/users", func() {
		m.Get("/", api.GetCurrentUser)
	})

	m.Group("/contacts", func() {
		m.Get("/", api.GetContacts)
	})

	m.Group("/labels", func() {
		m.Get("/", api.GetLabels)
	})

	m.Group("/messages", func() {
		m.Get("/count", api.GetMessagesCount)
	})

	m.Group("/conversations", func() {
		m.Get("/", api.GetConversations)
		m.Get("/count", api.GetConversationsCount)
	})

	m.Group("/events", func() {
		m.Get("/:event", api.GetEvent)
	})

	m.Post("/bugs/crash", binding.Json(CrashReq{}), api.Crash)
}
