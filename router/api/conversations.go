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
	filter := getMessagesFilter(ctx)

	conversations, total, err := api.backend.ListConversations(userId, filter)
	if err != nil {
		return
	}

	ctx.JSON(200, &ConversationsListResp{
		Resp: Resp{Ok},
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
		Resp: Resp{Ok},
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
		Resp: Resp{Ok},
		Conversation: conv,
		Messages: msgs,
	})
	return
}

func (api *Api) setConversationsRead(ctx *macaron.Context, req BatchReq, value int) {
	userId := api.getUserId(ctx)

	var respItems []*BatchRespItem

	for _, id := range req.IDs {
		r := &BatchRespItem{ ID: id }

		msgs, err := api.backend.ListConversationMessages(userId, id)
		if err != nil {
			r.Response = &ErrorResp{
				Resp: Resp{InternalServerError},
				ErrorDescription: err.Error(),
			}
		} else {
			for _, msg := range msgs {
				msg.IsRead = value

				_, err = api.backend.UpdateMessage(userId, &backend.MessageUpdate{
					Message: msg,
					IsRead: true,
				})

				if err != nil {
					r.Response = newErrorResp(err)
					break
				}
			}
		}

		if r.Response == nil {
			r.Response = &Resp{Ok}
		}

		respItems = append(respItems, r)
	}

	ctx.JSON(200, &BatchResp{
		Resp: Resp{Batch},
		Responses: respItems,
	})
}

func (api *Api) SetConversationsRead(ctx *macaron.Context, req BatchReq) {
	api.setConversationsRead(ctx, req, 1)
}

func (api *Api) SetConversationsUnread(ctx *macaron.Context, req BatchReq) {
	api.setConversationsRead(ctx, req, 0)
}
