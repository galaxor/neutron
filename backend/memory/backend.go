package memory

import (
	"github.com/emersion/neutron/backend"
)

type Backend struct {
	backend.DomainsBackend
	backend.ContactsBackend
	backend.LabelsBackend
	backend.ConversationsBackend
	backend.SendBackend
	backend.EventsBackend

	users map[string]*user
}

type user struct {
	*backend.User
	password string
}

func New() backend.Backend {
	bkd := &Backend{
		users: map[string]*user{},
	}

	bkd.DomainsBackend = NewDomainsBackend()
	bkd.ContactsBackend = NewContactsBackend()
	bkd.LabelsBackend = NewLabelsBackend()
	bkd.ConversationsBackend = NewConversationsBackend()
	bkd.SendBackend = NewEchoSendBackend(bkd.ConversationsBackend)
	bkd.EventsBackend = NewEventsBackend()

	return bkd
}
