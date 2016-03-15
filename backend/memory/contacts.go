package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListContacts(user string) (contacts []*backend.Contact, err error) {
	contacts = b.data[user].contacts
	return
}
