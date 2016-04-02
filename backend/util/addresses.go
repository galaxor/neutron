package util

import (
	"github.com/emersion/neutron/backend"
)

type EventedAddresses struct {
	backend.AddressesBackend
	events backend.EventsBackend
}

func (b *EventedAddresses) InsertAddress(user string, addr *backend.Address) (inserted *backend.Address, err error) {
	inserted, err = b.AddressesBackend.InsertAddress(user, addr)
	if err != nil {
		return
	}

	event := backend.NewUserEvent(&backend.User{ID: user})
	b.events.InsertEvent(user, event)
	return
}

func (b *EventedAddresses) DeleteAddress(user, id string) (err error) {
	err = b.AddressesBackend.DeleteAddress(user, id)
	if err != nil {
		return
	}

	event := backend.NewUserEvent(&backend.User{ID: user})
	b.events.InsertEvent(user, event)
	return
}

func NewEventedAddresses(addrs backend.AddressesBackend, events backend.EventsBackend) backend.AddressesBackend {
	return &EventedAddresses{
		AddressesBackend: addrs,
		events: events,
	}
}
