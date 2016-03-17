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
	InsertContact(user string, contact *Contact) (*Contact, error)
	UpdateContact(user string, update *ContactUpdate) (*Contact, error)
	DeleteContact(user, id string) error
	DeleteAllContacts(user string) error

	ListLabels(user string) ([]*Label, error)
	//InsertLabel(user string, label *Label) error
	//UpdateLabel(user string, update *LabelUpdate) error
	//DeleteLabel(user, id string) error

	GetMessage(user, id string) (*Message, error)
	ListMessages(user string, filter *MessagesFilter) ([]*Message, int, error)
	ListConversationMessages(user, id string) (msgs []*Message, err error)
	InsertMessage(user string, msg *Message) (*Message, error)
	UpdateMessage(user string, update *MessageUpdate) (*Message, error)
	DeleteMessage(user, id string) error

	GetConversation(user, id string) (conv *Conversation, err error)
	ListConversations(user string, filter *MessagesFilter) ([]*Conversation, int, error)
	CountConversations(user string) ([]*ConversationsCount, error)
	//DeleteConversation(user, id string) error
}
