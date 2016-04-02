package memory

import (
	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/util"
)

func Use(bkd *backend.Backend) {
	events := NewEvents()
	contacts := util.NewEventedContacts(NewContacts(), events)
	labels := util.NewEventedLabels(NewLabels(), events)
	attachments := NewAttachments()
	messages := NewMessages(attachments.(*Attachments))
	conversations := util.NewEventedConversations(NewConversations(messages.(*Messages)), events)
	send := util.NewEchoSend(conversations)
	domains := NewDomains()
	users := NewUsers()
	keys := NewKeys()

	bkd.Set(contacts, labels, conversations, send, domains, events, users, attachments, keys)
}
