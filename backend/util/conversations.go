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
