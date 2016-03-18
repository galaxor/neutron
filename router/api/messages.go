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

type SendMessageReq struct {
	Req
	ID string `json:"id"`
	Packages []*backend.MessagePackage
	AttachmentKeys []string // TODO
	ClearBody string
}

type SendMessageResp struct {
	Resp
	Sent *backend.Message
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
	msg.Type = backend.DraftType

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
	msgId := ctx.Params("id")

	msg := req.getMessage()
	msg.ID = msgId
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

func (api *Api) SendMessage(ctx *macaron.Context, req SendMessageReq) (err error) {
	userId := api.getUserId(ctx)
	msgId := ctx.Params("id")

	// Send each package
	for _, pkg := range req.Packages {
		err = api.backend.SendMessagePackage(userId, pkg)
		if err != nil {
			return
		}
	}

	// If clear body is available, send it to recipients without package
	if req.ClearBody != "" {
		var msg *backend.Message
		msg, err = api.backend.GetMessage(userId, msgId)
		if err != nil {
			return
		}

		for _, email := range msg.ToList {
			alreadySent := false
			for _, pkg := range req.Packages {
				if pkg.Address == email.Address {
					alreadySent = true
					break
				}
			}
			if alreadySent {
				continue
			}

			pkg := &backend.MessagePackage{
				Address: email.Address,
				Body: req.ClearBody,
			}

			err = api.backend.SendMessagePackage(userId, pkg)
			if err != nil {
				return
			}
		}
	}

	// Move message to Sent folder
	msg, err := api.backend.UpdateMessage(userId, &backend.MessageUpdate{
		Message: &backend.Message{
			ID: msgId,
			LabelIDs: []string{backend.SentLabel},
			Type: backend.SentType,
		},
		Type: true,
		LabelIDs: backend.ReplaceLabels,
	})
	if err != nil {
		return
	}

	ctx.JSON(200, &SendMessageResp{
		Resp: Resp{Ok},
		Sent: msg,
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
