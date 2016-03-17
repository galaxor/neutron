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
	Time int
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

type ConversationsCount struct {
	LabelID string
	Total int
	Unread int
}

type ConversationsFilter struct {
	Limit int
	Page int
	Label string
	Keyword string
	Address string // Address ID
	Attachments bool
	From string
	To string
	Begin int // Timestamp
	End int // Timestamp
	Sort string
	Desc bool
}
