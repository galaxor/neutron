package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListConversations(user, label string, limit, page int) (convs []*backend.Conversation, total int, err error) {
	// TODO: limit, page support	

	convs = b.data[user].conversations

	total = len(convs)

	return
}
