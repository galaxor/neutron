package api

import (
	"encoding/json"
	"log"

	"gopkg.in/macaron.v1"
)

type ClientType int

const (
	ClientEmail ClientType = 1
	ClientVPN
)

type CrashReq struct {
	Req
	OS string
	OSVersion string
	Browser string
	BrowserVersion string
	Client string
	ClientVersion string
	ClientType ClientType
	Debug json.RawMessage
}

func (api *Api) Crash(ctx *macaron.Context, req CrashReq) {
	log.Println("Client crashed:")
	log.Println(req)

	ctx.JSON(200, &Resp{Ok})
}
