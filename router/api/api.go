package api

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
)

type Req struct {}

type Resp struct {
	Code int
}

type ErrorResp struct {
	Resp
	Error string
	ErrorDescription string
}

func RegisterRoutes(m *macaron.Macaron) {
	m.Group("/auth", func() {
		m.Post("/", binding.Json(AuthReq{}), Auth)
		m.Post("/cookies", binding.Json(AuthCookiesReq{}), AuthCookies)
	})

	m.Group("/users", func() {
		m.Get("/", GetCurrentUser)
	})

	m.Group("/contacts", func() {
		m.Get("/", GetContacts)
	})

	m.Group("/labels", func() {
		m.Get("/", GetLabels)
	})

	m.Group("/messages", func() {
		m.Get("/count", GetMessagesCount)
	})

	m.Group("/conversations", func() {
		m.Get("/", GetConversations)
		m.Get("/count", GetConversationsCount)
	})

	m.Group("/events", func() {
		m.Get("/:event", GetEvent)
	})

	m.Post("/bugs/crash", binding.Json(CrashReq{}), Crash)
}
