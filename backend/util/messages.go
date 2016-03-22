package util

import (
	"time"

	"github.com/emersion/neutron/backend"
)

type EventedMessagesBackend struct {
	backend.MessagesBackend
	events backend.EventsBackend
}

func (b *EventedMessagesBackend) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	msg, err := b.MessagesBackend.InsertMessage(user, msg)

	if err == nil {
		event := backend.NewMessageDeltaEvent(msg.ID, backend.EventCreate, msg)
		b.events.InsertEvent(user, event)

		// TODO: add MessageCounts to event
	}

	return msg, err
}

func (b *EventedMessagesBackend) UpdateMessage(user string, update *backend.MessageUpdate) (*backend.Message, error) {
	msg, err := b.MessagesBackend.UpdateMessage(user, update)

	if err == nil {
		event := backend.NewMessageDeltaEvent(msg.ID, backend.EventUpdate, msg)
		b.events.InsertEvent(user, event)
	}

	return msg, err
}

func (b *EventedMessagesBackend) DeleteMessage(user, id string) error {
	err := b.MessagesBackend.DeleteMessage(user, id)

	if err == nil {
		event := backend.NewMessageDeltaEvent(id, backend.EventUpdate, nil)
		b.events.InsertEvent(user, event)

		// TODO: add MessageCounts to event
	}

	return err
}

func NewEventedMessagesBackend(bkd backend.MessagesBackend, events backend.EventsBackend) backend.MessagesBackend {
	return &EventedMessagesBackend{
		MessagesBackend: bkd,
		events: events,
	}
}


// A SendBackend that does nothing.
type NoopSendBackend struct {}

func (b *NoopSendBackend) SendMessagePackage(user string, msg *backend.OutgoingMessage) error {
	return nil // Do nothing
}

func NewNoopSendBackend() backend.SendBackend {
	return &NoopSendBackend{}
}


// A SendBackend that forwards all sent messages to a MessagesBackend.
type EchoSendBackend struct {
	target backend.MessagesBackend
}

func (b *EchoSendBackend) SendMessagePackage(user string, msg *backend.OutgoingMessage) error {
	newMsg := *msg.Message // Copy msg
	newMsg.Subject = "[EchoSendBackend forwarded message] " + newMsg.Subject
	newMsg.Body = msg.MessagePackage.Body
	newMsg.Time = time.Now().Unix()
	newMsg.LabelIDs = []string{backend.InboxLabel}
	newMsg.Type = 0
	newMsg.IsRead = 0

	_, err := b.target.InsertMessage(user, &newMsg)
	return err
}

func NewEchoSendBackend(target backend.MessagesBackend) backend.SendBackend {
	return &EchoSendBackend{
		target: target,
	}
}
