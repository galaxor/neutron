package api

import (
	"gopkg.in/macaron.v1"
)

func GetMessagesCount(ctx *macaron.Context) {
	GetConversationsCount(ctx) // TODO?
}
