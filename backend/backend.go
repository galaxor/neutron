// Contains a generic interface for backends.
package backend

// A backend takes care of storing all mailbox data.
type Backend interface {
	// Check if a username is available.
	IsUsernameAvailable(username string) (bool, error)
	// Get a user.
	GetUser(id string) (*User, error)
	// Check if the provided username and password are correct
	Auth(username, password string) (*User, error)
	// Insert a new user. Returns the newly created user.
	InsertUser(user *User, password string) (*User, error)
	// Update an existing user.
	UpdateUser(update *UserUpdate) error
	//UpdateUserPassword(id, password string) error
	//DeleteUser(id string) error

	// Get a public key for a user.
	GetPublicKey(email string) (string, error)

	// List all user's contacts.
	ListContacts(user string) ([]*Contact, error)
	// Insert a new contact.
	InsertContact(user string, contact *Contact) (*Contact, error)
	// Update an existing contact.
	UpdateContact(user string, update *ContactUpdate) (*Contact, error)
	// Delete a contact.
	DeleteContact(user, id string) error
	// Delete all contacts of a specific user.
	DeleteAllContacts(user string) error

	// List all user's labels.
	ListLabels(user string) ([]*Label, error)
	// Insert a new label.
	InsertLabel(user string, label *Label) (*Label, error)
	// Update an existing label.
	UpdateLabel(user string, update *LabelUpdate) (*Label, error)
	// Delete a label.
	DeleteLabel(user, id string) error

	// Get a message.
	GetMessage(user, id string) (*Message, error)
	// List all user's messages. A message filter can be provided.
	ListMessages(user string, filter *MessagesFilter) ([]*Message, int, error)
	// List all messages belonging to a conversation.
	ListConversationMessages(user, id string) (msgs []*Message, err error)
	// Count all user's messages by label.
	CountMessages(user string) ([]*MessagesCount, error)
	// Insert a new message.
	InsertMessage(user string, msg *Message) (*Message, error)
	// Update an existing message.
	UpdateMessage(user string, update *MessageUpdate) (*Message, error)
	// Permanently delete a message.
	DeleteMessage(user, id string) error
	// Send a message to an e-mail address.
	SendMessagePackage(user string, pkg *MessagePackage) error

	// Get a specific conversation.
	GetConversation(user, id string) (conv *Conversation, err error)
	// List all user's conversations. A message filter can be provided.
	ListConversations(user string, filter *MessagesFilter) ([]*Conversation, int, error)
	// Count all user's conversations by label.
	CountConversations(user string) ([]*MessagesCount, error)
	// Permanently delete a conversation.
	DeleteConversation(user, id string) error
}
