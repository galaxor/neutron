package backend

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

type Email struct {
	Address string
	Name string
}

type ConversationLabel struct {
	ID string
	Count int
	NumMessages int
	NumUnread int
}
