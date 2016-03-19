package util

import (
	"github.com/emersion/neutron/backend"
)

type EventedContactsBackend struct {
	backend.ContactsBackend
	events backend.EventsBackend
}

func (b *EventedContactsBackend) InsertContact(user string, contact *backend.Contact) (*backend.Contact, error) {
	contact, err := b.ContactsBackend.InsertContact(user, contact)

	if err == nil {
		event := backend.NewContactDeltaEvent(contact.ID, backend.EventCreate, contact)
		b.events.InsertEvent(user, event)
	}

	return contact, err
}

func (b *EventedContactsBackend) UpdateContact(user string, update *backend.ContactUpdate) (*backend.Contact, error) {
	contact, err := b.ContactsBackend.UpdateContact(user, update)

	if err == nil {
		event := backend.NewContactDeltaEvent(contact.ID, backend.EventUpdate, contact)
		b.events.InsertEvent(user, event)
	}

	return contact, err
}

func (b *EventedContactsBackend) DeleteContact(user, id string) error {
	err := b.ContactsBackend.DeleteContact(user, id)

	if err == nil {
		event := backend.NewContactDeltaEvent(id, backend.EventDelete, nil)
		b.events.InsertEvent(user, event)
	}

	return err
}

func NewEventedContactsBackend(bkd backend.ContactsBackend, events backend.EventsBackend) backend.ContactsBackend {
	return &EventedContactsBackend{
		ContactsBackend: bkd,
		events: events,
	}
}
