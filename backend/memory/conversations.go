package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) GetConversations(user, label string, limit, page int) (convs []*backend.Conversation, total int, err error) {
	convs = []*backend.Conversation{}

	total = len(convs)

	return
}
