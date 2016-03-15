package backend

type Backend interface {
	GetUser(id string) (*User, error)
	Auth(username, password string) (*User, error)

	ListContacts(user string) ([]*Contact, error)

	ListLabels(user string) ([]*Label, error)

	ListConversations(user, label string, limit, page int) ([]*Conversation, int, error)
	CountConversations(user string) ([]*ConversationsCount, error)
}
