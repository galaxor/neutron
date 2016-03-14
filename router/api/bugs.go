package api

import (
	"log"

	"gopkg.in/macaron.v1"
)

type CrashReq struct {
	Req
	OSVersion string
	Client string
	ClientVersion string
	Debug struct {
		State string `json:"state"`
		Error struct {
			Message string `json:"message"`
			Stack string `json:"stack"`
		} `json:"error"`
	}
}

func Crash(ctx *macaron.Context, req CrashReq) {
	log.Println("Client crashed:")
	log.Println(req)

	ctx.Resp.WriteHeader(200)
}
