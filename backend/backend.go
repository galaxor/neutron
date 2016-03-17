package backend

type Backend interface {
	IsUsernameAvailable(username string) (bool, error)
	GetUser(id string) (*User, error)
	Auth(username, password string) (*User, error)
	InsertUser(user *User, password string) (*User, error)
	UpdateUser(update *UserUpdate) error
	//UpdateUserPassword(id, password string) error
	//DeleteUser(id string) error

	ListContacts(user string) ([]*Contact, error)
	//InsertContact(user string, contact *Contact) error
	//UpdateContact(user string, update *ContactUpdate) error
	//DeleteContact(user, id string) error

	ListLabels(user string) ([]*Label, error)
	//InsertLabel(user string, label *Label) error
	//UpdateLabel(user string, update *LabelUpdate) error
	//DeleteLabel(user, id string) error

	GetMessage(user, id string) (*Message, error)
	UpdateMessage(user string, update *MessageUpdate) error
	//DeleteMessage(user, id string) error

	ListConversations(user string, filter *ConversationsFilter) ([]*Conversation, int, error)
	CountConversations(user string) ([]*ConversationsCount, error)
	GetConversation(user, id string) (conv *Conversation, err error)
	ListConversationMessages(user, id string) (msgs []*Message, err error)
}
