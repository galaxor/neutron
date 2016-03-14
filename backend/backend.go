package backend

type Backend interface {
	GetUser(id string) (*User, error)
	Auth(username, password string) (*User, error)

	GetContacts(user string) ([]*Contact, error)

	GetLabels(user string) ([]*Label, error)

	GetConversations(user, label string, limit, page int) ([]*Conversation, int, error)
}
