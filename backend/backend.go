// Contains a generic interface for backends.
package backend

// A backend takes care of storing all mailbox data.
type Backend struct {
	ContactsBackend
	LabelsBackend
	ConversationsBackend
	SendBackend
	DomainsBackend
	EventsBackend
	UsersBackend
	KeysBackend
}

// Set one or some of this backend's components.
func (b *Backend) Set(backends ...interface{}) {
	for _, bkd := range backends {
		if contacts, ok := bkd.(ContactsBackend); ok {
			b.ContactsBackend = contacts
		}
		if labels, ok := bkd.(LabelsBackend); ok {
			b.LabelsBackend = labels
		}
		if conversations, ok := bkd.(ConversationsBackend); ok {
			b.ConversationsBackend = conversations
		}
		if send, ok := bkd.(SendBackend); ok {
			b.SendBackend = send
		}
		if domains, ok := bkd.(DomainsBackend); ok {
			b.DomainsBackend = domains
		}
		if events, ok := bkd.(EventsBackend); ok {
			b.EventsBackend = events
		}
		if users, ok := bkd.(UsersBackend); ok {
			b.UsersBackend = users
		}
		if keys, ok := bkd.(KeysBackend); ok {
			b.KeysBackend = keys
		}
	}
}

func New() *Backend {
	return &Backend{}
}
