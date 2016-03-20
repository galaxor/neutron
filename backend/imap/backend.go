package imap

import (
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
}

func New() backend.Backend {
	bkd := &Backend{
		users: map[string]*backend.User{},
		connBackend: newConnBackend(),
	}

	messages := newMessagesBackend(bkd.connBackend)
	conversations := util.NewDummyConversationsBackend(messages)

	// TODO: do not use memory backends
	bkd.EventsBackend = memory.NewEventsBackend()
	bkd.DomainsBackend = memory.NewDomainsBackend()
	bkd.ContactsBackend = util.NewEventedContactsBackend(memory.NewContactsBackend(), bkd.EventsBackend)
	bkd.LabelsBackend = util.NewEventedLabelsBackend(memory.NewLabelsBackend(), bkd.EventsBackend)
	bkd.ConversationsBackend = util.NewEventedConversationsBackend(conversations, bkd.EventsBackend)
	bkd.SendBackend = util.NewEchoSendBackend(bkd.ConversationsBackend)

	return bkd
}
