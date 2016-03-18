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

	ctx.JSON(200, &MessagesCountResp{
		Resp: Resp{Ok},
		Counts: counts,
	})
	return
}

type ConversationResp struct {
	Resp
	Conversation *backend.Conversation
	Messages []*backend.Message
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

func (api *Api) batchUpdateConversationMessages(ctx *macaron.Context, req BatchReq, updater batchMessageUpdater) {
	userId := api.getUserId(ctx)

	var respItems []*BatchRespItem

	for _, id := range req.IDs {
		r := &BatchRespItem{ ID: id }
		respItems = append(respItems, r)

		msgs, err := api.backend.ListConversationMessages(userId, id)
		if err != nil {
			r.Response = &ErrorResp{
				Resp: Resp{InternalServerError},
				ErrorDescription: err.Error(),
			}
			continue
		}

		for _, msg := range msgs {
			// Create a new Message struct to prevent modifications on msg
			update := &backend.MessageUpdate{
				Message: &backend.Message{ ID: msg.ID },
			}
			updater(update)

			_, err = api.backend.UpdateMessage(userId, update)

			if err != nil {
				r.Response = newErrorResp(err)
				break
			}
		}

		if r.Response == nil {
			r.Response = &Resp{Ok}
		}
	}

	ctx.JSON(200, &BatchResp{
		Resp: Resp{Batch},
		Responses: respItems,
	})
}

func (api *Api) SetConversationsRead(ctx *macaron.Context, req BatchReq) {
	api.batchUpdateConversationMessages(ctx, req, func(update *backend.MessageUpdate) {
		update.IsRead = true
		update.Message.IsRead = 1
	})
}

func (api *Api) SetConversationsUnread(ctx *macaron.Context, req BatchReq) {
	api.batchUpdateConversationMessages(ctx, req, func(update *backend.MessageUpdate) {
		update.IsRead = true
		update.Message.IsRead = 0
	})
}

func (api *Api) SetConversationsStar(ctx *macaron.Context, req BatchReq) {
	api.batchUpdateConversationMessages(ctx, req, func(update *backend.MessageUpdate) {
		update.Starred = true
		update.LabelIDs = backend.AddLabels
		update.Message.LabelIDs = []string{backend.StarredLabel}
		update.Message.Starred = 1
	})
}

func (api *Api) SetConversationsUnstar(ctx *macaron.Context, req BatchReq) {
	api.batchUpdateConversationMessages(ctx, req, func(update *backend.MessageUpdate) {
		update.Starred = true
		update.LabelIDs = backend.RemoveLabels
		update.Message.LabelIDs = []string{backend.StarredLabel}
		update.Message.Starred = 0
	})
}
