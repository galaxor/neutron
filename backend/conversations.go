package backend

type Conversation struct {
	ID string
	Order int
	Subject string
	Senders []*ConversationAddress
	Recipients []*ConversationAddress
	NumMessages int
	NumUnread int
	NumAttachments int
	ExpirationTime int
	TotalSize int
	Time int
	LabelIDs []string
	Labels []*ConversationLabel
}

type ConversationAddress struct {
	Address string
	Name string
}

type ConversationLabel struct {
	ID string
	Count int
	NumMessages int
	NumUnread int
}

func GetConversations(user, label string, limit, page int) (conversations []*Conversation, total int, err error) {
	conversations = []*Conversation{}

	total = len(conversations)

	return
}
