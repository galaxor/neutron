package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type ConversationsListResp struct {
	Resp
	Total int
	Conversations []*backend.Conversation
}

type ConversationsCountResp struct {
	Resp
	Counts []*backend.ConversationsCount
}

type ConversationResp struct {
	Resp
	Conversation *backend.Conversation
	Messages []*backend.Message
}

func (api *Api) ListConversations(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)

	filter := &backend.ConversationsFilter{
		Limit: ctx.QueryInt("Limit"),
		Page: ctx.QueryInt("Page"),
		Label: ctx.Query("Label"),
		Keyword: ctx.Query("Keyword"),
		Address: ctx.Query("Address"),
		Attachments: (ctx.Query("Attachments") == "1"),
		From: ctx.Query("From"),
		To: ctx.Query("To"),
		Begin: ctx.QueryInt("Begin"),
		End: ctx.QueryInt("End"),
		Sort: ctx.Query("Sort"),
		Desc: (ctx.Query("Desc") == "1"),
	}

	conversations, total, err := api.backend.ListConversations(userId, filter)
	if err != nil {
		return
	}

	ctx.JSON(200, &ConversationsListResp{
		Resp: Resp{1000},
		Total: total,
		Conversations: conversations,
	})
	return
}

func (api *Api) GetConversationsCount(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)

	counts, err := api.backend.CountConversations(userId)
	if err != nil {
		return
	}

	ctx.JSON(200, &ConversationsCountResp{
		Resp: Resp{1000},
		Counts: counts,
	})
	return
}

func (api *Api) GetConversation(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)
	convId := ctx.Params("id")

	conv, err := api.backend.GetConversation(userId, convId)
	if err != nil {
		return
	}

	msgs, err := api.backend.ListConversationMessages(userId, convId)
	if err != nil {
		return
	}

	ctx.JSON(200, &ConversationResp{
		Resp: Resp{1000},
		Conversation: conv,
		Messages: msgs,
	})
	return
}
