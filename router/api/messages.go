package api

import (
	"time"

	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type MessagesListResp struct {
	Resp
	Total int
	Messages []*backend.Message
}

type MessageReq struct {
	Req
	Message *backend.Message
	ID string `json:"id"`
}

func (req MessageReq) getMessage() *backend.Message {
	return &backend.Message{
		ID: req.ID,
		ToList: req.Message.ToList,
		CCList: req.Message.CCList,
		BCCList: req.Message.BCCList,
		Subject: req.Message.Subject,
		IsRead: req.Message.IsRead,
		AddressID: req.Message.AddressID,
		Body: req.Message.Body,
	}
}

type MessageResp struct {
	Resp
	Message *backend.Message
}

func getMessagesFilter(ctx *macaron.Context) *backend.MessagesFilter {
	return &backend.MessagesFilter{
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
}

func (api *Api) ListMessages(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)
	filter := getMessagesFilter(ctx)

	msgs, total, err := api.backend.ListMessages(userId, filter)
	if err != nil {
		return
	}

	ctx.JSON(200, &MessagesListResp{
		Resp: Resp{Ok},
		Total: total,
		Messages: msgs,
	})
	return
}

func (api *Api) GetMessagesCount(ctx *macaron.Context) {
	api.GetConversationsCount(ctx) // TODO?
}

func (api *Api) setMessagesRead(ctx *macaron.Context, req BatchReq, value int) {
	userId := api.getUserId(ctx)

	var respItems []*BatchRespItem

	for _, id := range req.IDs {
		_, err := api.backend.UpdateMessage(userId, &backend.MessageUpdate{
			Message: &backend.Message{
				ID: id,
				IsRead: 1,
			},
			IsRead: true,
		})

		r := &BatchRespItem{ ID: id }
		if err != nil {
			r.Response = newErrorResp(err)
		} else {
			r.Response = &Resp{Ok}
		}
		respItems = append(respItems, r)
	}

	ctx.JSON(200, newBatchResp(respItems))
}

func (api *Api) SetMessagesRead(ctx *macaron.Context, req BatchReq) {
	api.setMessagesRead(ctx, req, 1)
}

func (api *Api) SetMessagesUnread(ctx *macaron.Context, req BatchReq) {
	api.setMessagesRead(ctx, req, 0)
}

func (api *Api) CreateDraft(ctx *macaron.Context, req MessageReq) (err error) {
	userId := api.getUserId(ctx)

	user, err := api.backend.GetUser(userId)
	if err != nil {
		return
	}

	msg := req.getMessage()
	msg.Attachments = []*backend.Attachment{}
	msg.LabelIDs = []string{backend.DraftsLabel}
	msg.Time = time.Now().Unix()
	msg.Type = 1

	for _, address := range user.Addresses {
		if address.ID == msg.AddressID {
			msg.Sender = address.GetEmail()
			break
		}
	}

	msg, err = api.backend.InsertMessage(userId, msg)
	if err != nil {
		return
	}

	ctx.JSON(200, &MessageResp{
		Resp: Resp{Ok},
		Message: msg,
	})
	return
}

func (api *Api) UpdateDraft(ctx *macaron.Context, req MessageReq) (err error) {
	userId := api.getUserId(ctx)

	msg := req.getMessage()
	msg.Time = time.Now().Unix()

	msg, err = api.backend.UpdateMessage(userId, &backend.MessageUpdate{
		Message: msg,
		ToList: true,
		CCList: true,
		BCCList: true,
		Subject: true,
		IsRead: true,
		AddressID: true,
		Body: true,
		Time: true,
	})
	if err != nil {
		return
	}

	ctx.JSON(200, &MessageResp{
		Resp: Resp{Ok},
		Message: msg,
	})
	return
}

func (api *Api) DeleteMessages(ctx *macaron.Context, req BatchReq) {
	userId := api.getUserId(ctx)

	var respItems []*BatchRespItem

	for _, id := range req.IDs {
		r := &BatchRespItem{ ID: id }
		respItems = append(respItems, r)

		err := api.backend.DeleteMessage(userId, id)
		if err != nil {
			r.Response = newErrorResp(err)
		} else {
			r.Response = &Resp{Ok}
		}
	}

	ctx.JSON(200, newBatchResp(respItems))
}
