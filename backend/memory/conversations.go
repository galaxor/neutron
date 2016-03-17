package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) populateConversation(user string, conv *backend.Conversation) error {
	msgs, err := b.ListConversationMessages(user, conv.ID)
	if err != nil {
		return err
	}

	conv.NumMessages = 0
	conv.NumUnread = 0
	conv.Labels = nil
	conv.LabelIDs = nil

	for _, msg := range msgs {
		conv.NumMessages++
		if msg.IsRead == 0 {
			conv.NumUnread++
		}

		for _, labelId := range msg.LabelIDs {
			var label *backend.ConversationLabel
			for _, l := range conv.Labels {
				if l.ID == labelId {
					label = l
					break
				}
			}

			if label == nil {
				label = &backend.ConversationLabel{ ID: labelId }
				conv.Labels = append(conv.Labels, label)
				conv.LabelIDs = append(conv.LabelIDs, labelId)
			}

			label.NumMessages++
			if msg.IsRead == 0 {
				label.NumUnread++
			}
		}
	}

	return nil
}

func (b *Backend) ListConversations(user string, filter *backend.ConversationsFilter) (convs []*backend.Conversation, total int, err error) {
	// TODO: filter according to label

	all := b.data[user].conversations
	filtered := []*backend.Conversation{}

	for _, c := range all {
		b.populateConversation(user, c)

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

		convs = filtered[from:to]
	} else {
		convs = filtered
	}

	return
}

func (b *Backend) CountConversations(user string) (counts []*backend.ConversationsCount, err error) {
	convs := b.data[user].conversations

	indexes := map[string]int{}

	for _, c := range convs {
		for _, label := range c.LabelIDs {
			var count *backend.ConversationsCount
			if i, ok := indexes[label]; ok {
				count = counts[i]
			} else {
				indexes[label] = len(counts)
				count = &backend.ConversationsCount{ LabelID: label }
			}

			count.Total++
			if c.NumUnread > 0 {
				count.Unread++
			}
		}
	}

	return
}

func (b *Backend) GetConversation(user, id string) (conv *backend.Conversation, err error) {
	for _, c := range b.data[user].conversations {
		if c.ID == id {
			conv = c
			b.populateConversation(user, conv)
			return
		}
	}

	err = errors.New("No such conversation")
	return
}

func (b *Backend) ListConversationMessages(user, id string) (msgs []*backend.Message, err error) {
	for _, m := range b.data[user].messages {
		if m.ConversationID == id {
			msgs = append(msgs, m)
		}
	}
	return
}
