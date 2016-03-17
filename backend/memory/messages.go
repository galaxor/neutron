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

func (b *Backend) ListMessages(user string, filter *backend.MessagesFilter) (msgs []*backend.Message, total int, err error) {
	all := b.data[user].messages
	filtered := []*backend.Message{}

	for _, c := range all {
		if filter.Label != "" {
			matches := false
			for _, lbl := range c.LabelIDs {
				if lbl == filter.Label {
					matches = true
					break
				}
			}

			if !matches {
				continue
			}
		}

		// TODO: other filter fields support

		filtered = append(filtered, c)
	}

	total = len(filtered)

	if filter.Limit > 0 && filter.Page >= 0 {
		from := filter.Limit * filter.Page
		to := filter.Limit * (filter.Page + 1)
		if from < 0 {
			from = 0
		}
		if to > total {
			to = total
		}

		msgs = filtered[from:to]
	} else {
		msgs = filtered
	}

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

func (b *Backend) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	msg.ID = generateId()
	if msg.ConversationID == "" {
		msg.ConversationID = generateId()
	}

	b.data[user].messages = append(b.data[user].messages, msg)
	return msg, nil
}

func (b *Backend) UpdateMessage(user string, update *backend.MessageUpdate) (msg *backend.Message, err error) {
	updated := update.Message

	i, err := b.getMessageIndex(user, updated.ID)
	if err != nil {
		return
	}

	msg = b.data[user].messages[i]

	if update.ToList {
		msg.ToList = updated.ToList
	}
	if update.CCList {
		msg.CCList = updated.CCList
	}
	if update.BCCList {
		msg.BCCList = updated.BCCList
	}
	if update.Subject {
		msg.Subject = updated.Subject
	}
	if update.IsRead {
		msg.IsRead = updated.IsRead
	}
	if update.Type {
		msg.Type = updated.Type
	}
	if update.AddressID {
		msg.AddressID = updated.AddressID
	}
	if update.Body {
		msg.Body = updated.Body
	}
	if update.Time {
		msg.Time = updated.Time
	}

	if update.LabelIDs != backend.KeepLabels {
		switch update.LabelIDs {
		case backend.ReplaceLabels:
			msg.LabelIDs = updated.LabelIDs
		case backend.AddLabels:
			for _, lblToAdd := range updated.LabelIDs {
				found := false
				for _, lbl := range msg.LabelIDs {
					if lbl == lblToAdd {
						found = true
						break
					}
				}
				if !found {
					msg.LabelIDs = append(msg.LabelIDs, lblToAdd)
				}
			}
		case backend.RemoveLabels:
			var labels []string
			for _, lbl := range updated.LabelIDs {
				found := false
				for _, lblToRemove := range updated.LabelIDs {
					if lbl == lblToRemove {
						found = true
						break
					}
				}
				if !found {
					labels = append(labels, lbl)
				}
			}
			msg.LabelIDs = labels
		}
	}

	return
}

func (b *Backend) DeleteMessage(user, id string) error {
	i, err := b.getMessageIndex(user, id)
	if err != nil {
		return err
	}

	messages := b.data[user].messages
	b.data[user].messages = append(messages[:i], messages[i+1:]...)

	return nil
}

func (b *Backend) SendMessagePackage(user string, pkg *backend.MessagePackage) error {
	return nil // Do nothing
}
