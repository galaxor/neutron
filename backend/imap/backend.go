package imap

import (
	"errors"

	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/memory"
	"github.com/emersion/neutron/backend/util"
)

type Backend struct {
	backend.DomainsBackend
	backend.ContactsBackend
	backend.LabelsBackend
	backend.ConversationsBackend
	backend.SendBackend
	backend.EventsBackend

	*connBackend

	users map[string]*backend.User
	passwords map[string]string
}

func (b *Backend) Set(item interface{}) error {
	switch val := item.(type) {
	case backend.SendBackend:
		b.SendBackend = val
	default:
		return errors.New("Unsupported backend")
	}
	return nil
}

func New(config *Config) backend.Backend {
	bkd := &Backend{
		connBackend: newConnBackend(config),

		users: map[string]*backend.User{},
		passwords: map[string]string{},
	}

	messages := newMessagesBackend(bkd.connBackend)
	conversations := util.NewDummyConversationsBackend(messages)

	// TODO: do not use memory backends
	bkd.EventsBackend = memory.NewEventsBackend()
	bkd.DomainsBackend = memory.NewDomainsBackend()
	bkd.ContactsBackend = util.NewEventedContactsBackend(memory.NewContactsBackend(), bkd.EventsBackend)
	bkd.LabelsBackend = util.NewEventedLabelsBackend(memory.NewLabelsBackend(), bkd.EventsBackend)
	bkd.ConversationsBackend = util.NewEventedConversationsBackend(conversations, bkd.EventsBackend)
	bkd.SendBackend = util.NewNoopSendBackend()

	return bkd
}
