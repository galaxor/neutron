package backend

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
