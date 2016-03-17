package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListContacts(user string) (contacts []*backend.Contact, err error) {
	contacts = b.data[user].contacts
	return
}

func (b *Backend) InsertContact(user string, contact *backend.Contact) (*backend.Contact, error) {
	contact.ID = generateId()
	b.data[user].contacts = append(b.data[user].contacts, contact)
	return contact, nil
}

func (b *Backend) getContactIndex(user, id string) (int, error) {
	for i, contact := range b.data[user].contacts {
		if contact.ID == id {
			return i, nil
		}
	}

	return -1, errors.New("No such contact")
}

func (b *Backend) UpdateContact(user string, update *backend.ContactUpdate) (*backend.Contact, error) {
	updated := update.Contact

	i, err := b.getContactIndex(user, updated.ID)
	if err != nil {
		return nil, err
	}

	contact := b.data[user].contacts[i]

	if update.Name {
		contact.Name = updated.Name
	}
	if update.Email {
		contact.Email = updated.Email
	}

	return contact, nil
}

func (b *Backend) DeleteContact(user, id string) error {
	i, err := b.getContactIndex(user, id)
	if err != nil {
		return err
	}

	contacts := b.data[user].contacts
	b.data[user].contacts = append(contacts[:i], contacts[i+1:]...)

	return nil
}

func (b *Backend) DeleteAllContacts(user string) error {
	b.data[user].contacts = nil
	return nil
}
