package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type ContactsBackend struct {
	contacts map[string][]*backend.Contact
}

func (b *ContactsBackend) ListContacts(user string) (contacts []*backend.Contact, err error) {
	contacts = b.contacts[user]
	return
}

func (b *ContactsBackend) InsertContact(user string, contact *backend.Contact) (*backend.Contact, error) {
	contact.ID = generateId()
	b.contacts[user] = append(b.contacts[user], contact)
	return contact, nil
}

func (b *ContactsBackend) getContactIndex(user, id string) (int, error) {
	for i, contact := range b.contacts[user] {
		if contact.ID == id {
			return i, nil
		}
	}

	return -1, errors.New("No such contact")
}

func (b *ContactsBackend) UpdateContact(user string, update *backend.ContactUpdate) (*backend.Contact, error) {
	updated := update.Contact

	i, err := b.getContactIndex(user, updated.ID)
	if err != nil {
		return nil, err
	}

	contact := b.contacts[user][i]

	if update.Name {
		contact.Name = updated.Name
	}
	if update.Email {
		contact.Email = updated.Email
	}

	return contact, nil
}

func (b *ContactsBackend) DeleteContact(user, id string) error {
	i, err := b.getContactIndex(user, id)
	if err != nil {
		return err
	}

	contacts := b.contacts[user]
	b.contacts[user] = append(contacts[:i], contacts[i+1:]...)

	return nil
}

func (b *ContactsBackend) DeleteAllContacts(user string) error {
	b.contacts[user] = nil
	return nil
}

func NewContactsBackend() backend.ContactsBackend {
	return &ContactsBackend{
		contacts: map[string][]*backend.Contact{},
	}
}
