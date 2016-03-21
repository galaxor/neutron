package backend

// Stores conversations data.
type ConversationsBackend interface {
	MessagesBackend

	// List all messages belonging to a conversation.
	ListConversationMessages(user, id string) (msgs []*Message, err error)

	// Get a specific conversation.
	GetConversation(user, id string) (conv *Conversation, err error)
	// List all user's conversations. A message filter can be provided.
	ListConversations(user string, filter *MessagesFilter) ([]*Conversation, int, error)
	// Count all user's conversations by label.
	CountConversations(user string) ([]*MessagesCount, error)
	// Permanently delete a conversation.
	DeleteConversation(user, id string) error
}

// A conversation is a sequence of messages.
type Conversation struct {
	ID string
	Order int
	Subject string
	Senders []*Email
	Recipients []*Email
	NumMessages int
	NumUnread int
	NumAttachments int
	ExpirationTime int
	TotalSize int
	Time int64
	LabelIDs []string
	Labels []*ConversationLabel
}

// An email contains an address and a name.
type Email struct {
	Address string
	Name string
}

// Contains messages counts by labels.
type ConversationLabel struct {
	ID string
	Count int
	NumMessages int
	NumUnread int
}
