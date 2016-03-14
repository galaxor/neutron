package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type ConversationsResp struct {
	Resp
	Total int
	Conversations []*backend.Conversation
}

type ConversationsCountResp struct {
	Resp
	Counts []*ConversationsCount
}

type ConversationsCount struct {
	LabelID string
	Total int
	Unread int
}

func GetConversations(ctx *macaron.Context) (err error) {
	label := ctx.Params("Label")
	limit := ctx.ParamsInt("Limit")
	page := ctx.ParamsInt("Page")

	conversations, total, err := backend.GetConversations(userId, label, limit, page)
	if err != nil {
		return
	}

	ctx.JSON(200, &ConversationsResp{
		Resp: Resp{1000},
		Total: total,
		Conversations: conversations,
	})
	return
}

func GetConversationsCount(ctx *macaron.Context) (err error) {
	// TODO
	ctx.JSON(200, &ConversationsCountResp{
		Resp: Resp{1000},
		Counts: []*ConversationsCount{
			&ConversationsCount{
				LabelID: "0",
				Total: 0,
				Unread: 0,
			},
		},
	})
	return
}
