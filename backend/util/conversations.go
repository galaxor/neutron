package util

import (
	"github.com/emersion/neutron/backend"
)

type EventedConversationsBackend struct {
	backend.ConversationsBackend
	messages backend.MessagesBackend
	events backend.EventsBackend
}

func (b *EventedConversationsBackend) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	hadConversation := (msg.ConversationID != "")

	msg, err := b.messages.InsertMessage(user, msg)

	if err == nil && msg.ConversationID != "" {
		action := backend.EventCreate
		if hadConversation {
			action = backend.EventUpdate
		}

		conv, err := b.GetConversation(user, msg.ConversationID)
		if err == nil {
			event := backend.NewConversationDeltaEvent(msg.ConversationID, action, conv)
			b.events.InsertEvent(user, event)
		}

		// TODO: add ConversationCounts to event
	}

	return msg, err
}

func (b *EventedConversationsBackend) UpdateMessage(user string, update *backend.MessageUpdate) (*backend.Message, error) {
	msg, err := b.messages.UpdateMessage(user, update)

	if err == nil && msg.ConversationID != "" {
		conv, err := b.GetConversation(user, msg.ConversationID)
		if err == nil {
			event := backend.NewConversationDeltaEvent(msg.ConversationID, backend.EventUpdate, conv)
			b.events.InsertEvent(user, event)
		}
	}

	return msg, err
}

func (b *EventedConversationsBackend) DeleteMessage(user, id string) error {
	msg, _ := b.GetMessage(user, id)

	err := b.messages.DeleteMessage(user, id)

	if err == nil && msg != nil && msg.ConversationID != "" {
		conv, _ := b.GetConversation(user, msg.ConversationID)

		action := backend.EventUpdate
		if conv == nil {
			action = backend.EventDelete
		}

		event := backend.NewConversationDeltaEvent(msg.ConversationID, action, conv)
		b.events.InsertEvent(user, event)

		// TODO: add ConversationCounts to event
	}

	return err
}

func NewEventedConversationsBackend(bkd backend.ConversationsBackend, events backend.EventsBackend) backend.ConversationsBackend {
	return &EventedConversationsBackend{
		ConversationsBackend: bkd,
		messages: NewEventedMessagesBackend(bkd, events),
		events: events,
	}
}


// A conversations backend that builds one conversation per message (no threads).
type DummyConversationsBackend struct {
	backend.MessagesBackend
}

func (b *DummyConversationsBackend) ListConversationMessages(user, id string) ([]*backend.Message, error) {
	msg, err := b.GetMessage(user, id)
	if err != nil {
		return nil, err
	}
	msg.ConversationID = id
	return []*backend.Message{msg}, nil
}

func (b *DummyConversationsBackend) buildConversation(msg *backend.Message) *backend.Conversation {
	conv := &backend.Conversation{
		ID: msg.ID,
		Order: msg.Order,
		NumMessages: 1,
		NumUnread: 1 - msg.IsRead,
		Time: msg.Time,
		Subject: msg.Subject,
		Senders: []*backend.Email{msg.Sender},
		Recipients: msg.ToList,
		TotalSize: msg.Size,
		LabelIDs: msg.LabelIDs,
	}

	for _, lbl := range msg.LabelIDs {
		conv.Labels = append(conv.Labels, &backend.ConversationLabel{
			ID: lbl,
			NumMessages: 1,
			NumUnread: 1 - msg.IsRead,
		})
	}

	return conv
}

func (b *DummyConversationsBackend) GetConversation(user, id string) (*backend.Conversation, error) {
	msg, err := b.GetMessage(user, id)
	if err != nil {
		return nil, err
	}
	return b.buildConversation(msg), nil
}

func (b *DummyConversationsBackend) ListConversations(user string, filter *backend.MessagesFilter) ([]*backend.Conversation, int, error) {
	msgs, total, err := b.ListMessages(user, filter)
	if err != nil {
		return nil, -1, err
	}

	convs := make([]*backend.Conversation, len(msgs))
	for i, msg := range msgs {
		convs[i] = b.buildConversation(msg)
	}

	return convs, total, nil
}

func (b *DummyConversationsBackend) CountConversations(user string) ([]*backend.MessagesCount, error) {
	return b.CountMessages(user)
}

func (b *DummyConversationsBackend) DeleteConversation(user, id string) error {
	return b.DeleteMessage(user, id)
}

func (b *DummyConversationsBackend) GetMessage(user, id string) (*backend.Message, error) {
	msg, err := b.MessagesBackend.GetMessage(user, id)

	if err == nil {
		msg.ConversationID = msg.ID
	}

	return msg, err
}

func (b *DummyConversationsBackend) ListMessages(user string, filter *backend.MessagesFilter) ([]*backend.Message, int, error) {
	msgs, total, err := b.MessagesBackend.ListMessages(user, filter)

	if err == nil {
		for _, msg := range msgs {
			msg.ConversationID = msg.ID
		}
	}

	return msgs, total, err
}

func (b *DummyConversationsBackend) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	msg, err := b.MessagesBackend.InsertMessage(user, msg)

	if err == nil {
		msg.ConversationID = msg.ID
	}

	return msg, err
}

func (b *DummyConversationsBackend) UpdateMessage(user string, update *backend.MessageUpdate) (*backend.Message, error) {
	msg, err := b.MessagesBackend.UpdateMessage(user, update)

	if err == nil {
		msg.ConversationID = msg.ID
	}

	return msg, err
}

func NewDummyConversationsBackend(messages backend.MessagesBackend) backend.ConversationsBackend {
	return &DummyConversationsBackend{messages}
}
