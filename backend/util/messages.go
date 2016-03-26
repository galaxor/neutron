package util

import (
	"github.com/emersion/neutron/backend"
)

type EventedMessages struct {
	backend.MessagesBackend
	events backend.EventsBackend
}

func (b *EventedMessages) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	msg, err := b.MessagesBackend.InsertMessage(user, msg)

	if err == nil {
		event := backend.NewMessageDeltaEvent(msg.ID, backend.EventCreate, msg)
		b.events.InsertEvent(user, event)

		// TODO: add MessageCounts to event
	}

	return msg, err
}

func (b *EventedMessages) UpdateMessage(user string, update *backend.MessageUpdate) (*backend.Message, error) {
	msg, err := b.MessagesBackend.UpdateMessage(user, update)

	if err == nil {
		event := backend.NewMessageDeltaEvent(msg.ID, backend.EventUpdate, msg)
		b.events.InsertEvent(user, event)
	}

	return msg, err
}

func (b *EventedMessages) DeleteMessage(user, id string) error {
	err := b.MessagesBackend.DeleteMessage(user, id)

	if err == nil {
		event := backend.NewMessageDeltaEvent(id, backend.EventUpdate, nil)
		b.events.InsertEvent(user, event)

		// TODO: add MessageCounts to event
	}

	return err
}

func NewEventedMessages(bkd backend.MessagesBackend, events backend.EventsBackend) backend.MessagesBackend {
	return &EventedMessages{
		MessagesBackend: bkd,
		events: events,
	}
}
