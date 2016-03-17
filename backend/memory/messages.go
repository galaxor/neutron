package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) getMessageIndex(user, id string) (int, error) {
	for i, m := range b.data[user].messages {
		if m.ID == id {
			return i, nil
		}
	}

	return -1, errors.New("No such message")
}

func (b *Backend) GetMessage(user, id string) (msg *backend.Message, err error) {
	i, err := b.getMessageIndex(user, id)
	if err != nil {
		return
	}

	msg = b.data[user].messages[i]
	return
}

func (b *Backend) SetMessageRead(user, id string, value int) (err error) {
	i, err := b.getMessageIndex(user, id)
	if err != nil {
		return
	}

	b.data[user].messages[i].IsRead = value
	return
}
