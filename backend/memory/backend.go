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

	users map[string]*user
}

type user struct {
	*backend.User
	password string
}

func New() backend.Backend {
	return &Backend{
		DomainsBackend: NewDomainsBackend(),
		ContactsBackend: NewContactsBackend(),
		LabelsBackend: NewLabelsBackend(),
		ConversationsBackend: NewConversationsBackend(),
		SendBackend: NewSendBackend(),
		users: map[string]*user{},
	}
}
