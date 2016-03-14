package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) GetContacts(user string) (contacts []*backend.Contact, err error) {
	contacts = []*backend.Contact{
		&backend.Contact{
			ID: "contact_id",
			Name: "Myself :)",
			Email: "neutron@example.org",
		},
	}
	return
}
