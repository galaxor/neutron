package api

import (
	"gopkg.in/macaron.v1"
)

func (api *Api) GetMessagesCount(ctx *macaron.Context) {
	api.GetConversationsCount(ctx) // TODO?
}
