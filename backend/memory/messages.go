package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func populateMessage(msg *backend.Message) {
	if msg.ToList == nil {
		msg.ToList = []*backend.Email{}
	}
	if msg.CCList == nil {
		msg.CCList = []*backend.Email{}
	}
	if msg.BCCList == nil {
		msg.BCCList = []*backend.Email{}
	}
	if msg.Attachments == nil {
		msg.Attachments = []*backend.Attachment{}
	}
	if msg.LabelIDs == nil {
		msg.LabelIDs = []string{}
	}

	if msg.Sender != nil {
		msg.SenderAddress = msg.Sender.Address
		msg.SenderName = msg.Sender.Name
	}

	if backend.IsEncrypted(msg.Body) {
		msg.IsEncrypted = 1
	}
}

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
	populateMessage(msg)
	return
}

func (b *Backend) ListMessages(user string, filter *backend.MessagesFilter) (msgs []*backend.Message, total int, err error) {
	all := b.data[user].messages
	filtered := []*backend.Message{}

	for _, msg := range all {
		if filter.Label != "" {
			matches := false
			for _, lbl := range msg.LabelIDs {
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

		populateMessage(msg)
		filtered = append(filtered, msg)
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
			populateMessage(msg)
			msgs = append(msgs, msg)
		}
	}
	return
}

func (b *Backend) CountMessages(user string) (counts []*backend.MessagesCount, err error) {
	indexes := map[string]int{}

	for _, msg := range b.data[user].messages {
		for _, label := range msg.LabelIDs {
			var count *backend.MessagesCount
			if i, ok := indexes[label]; ok {
				count = counts[i]
			} else {
				indexes[label] = len(counts)
				count = &backend.MessagesCount{ LabelID: label }
				counts = append(counts, count)
			}

			count.Total++
			if msg.IsRead == 0 {
				count.Unread++
			}
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
	populateMessage(msg)
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
			labels := []string{}
			for _, lbl := range msg.LabelIDs {
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

	populateMessage(msg)
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
