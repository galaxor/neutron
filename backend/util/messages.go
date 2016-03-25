package util

import (
	"time"

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


// A SendBackend that does nothing.
type NoopSend struct {}

func (b *NoopSend) SendMessagePackage(user string, msg *backend.OutgoingMessage) error {
	return nil // Do nothing
}

func NewNoopSend() backend.SendBackend {
	return &NoopSend{}
}


// A SendBackend that forwards all sent messages to a MessagesBackend.
type EchoSend struct {
	target backend.MessagesBackend
}

func (b *EchoSend) SendMessagePackage(user string, msg *backend.OutgoingMessage) error {
	newMsg := *msg.Message // Copy msg
	newMsg.Subject = "[EchoSend forwarded message] " + newMsg.Subject
	newMsg.Body = msg.MessagePackage.Body
	newMsg.Time = time.Now().Unix()
	newMsg.LabelIDs = []string{backend.InboxLabel}
	newMsg.Type = 0
	newMsg.IsRead = 0

	_, err := b.target.InsertMessage(user, &newMsg)
	return err
}

func NewEchoSend(target backend.MessagesBackend) backend.SendBackend {
	return &EchoSend{
		target: target,
	}
}
