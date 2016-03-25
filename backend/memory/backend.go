package memory

import (
	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/util"
)

func Use(bkd *backend.Backend) {
	events := NewEvents()
	contacts := util.NewEventedContacts(NewContacts(), events)
	labels := util.NewEventedLabels(NewLabels(), events)
	conversations := util.NewEventedConversations(NewConversations(), events)
	send := util.NewEchoSend(conversations)
	domains := NewDomains()
	users := NewUsers()

	bkd.Set(contacts, labels, conversations, send, domains, events, users)
}
