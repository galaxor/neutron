package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type Backend struct {
	backend.DomainsBackend

	data map[string]*userData
}

type userData struct {
	user *backend.User
	password string
	contacts []*backend.Contact
	messages []*backend.Message
	labels []*backend.Label
}

func (b *Backend) getUserData(id string) (*userData, error) {
	item, ok := b.data[id]
	if !ok {
		return nil, errors.New("No such user")
	}
	return item, nil
}

func New() backend.Backend {
	return &Backend{
		DomainsBackend: NewDomainsBackend(),
	}
}
