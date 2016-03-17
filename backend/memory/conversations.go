package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListConversations(user string, filter *backend.ConversationsFilter) (convs []*backend.Conversation, total int, err error) {
	// TODO: filter according to label

	all := b.data[user].conversations
	filtered := []*backend.Conversation{}

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
