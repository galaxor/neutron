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

func (b *Backend) ListConversationMessages(user, id string) (msgs []*backend.Message, err error) {
	for _, msg := range b.data[user].messages {
		if msg.ConversationID == id {
			msgs = append(msgs, msg)
		}
	}
	return
}

func (b *Backend) UpdateMessage(user string, update *backend.MessageUpdate) (err error) {
	updated := update.Message

	i, err := b.getMessageIndex(user, updated.ID)
	if err != nil {
		return
	}

	msg := b.data[user].messages[i]

	if update.IsRead {
		msg.IsRead = updated.IsRead
	}

	return
}
