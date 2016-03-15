package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListConversations(user, label string, limit, page int) (convs []*backend.Conversation, total int, err error) {
	allConvs := b.data[user].conversations
	total = len(allConvs)

	from := limit * page
	to := limit * (page + 1)
	if from < 0 {
		from = 0
	}
	if to > total {
		to = total
	}

	convs = allConvs[from:to]
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
