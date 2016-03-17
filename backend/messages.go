package backend

type Message struct {
	ID string
	Order int
	ConversationID string
	Subject string
	IsRead int
	Type int
	SenderAddress string
	SenderName string
	Sender *Email
	ToList []*Email
	CCList []*Email
	BCCList []*Email
	Time int
	Size int
	NumAttachments int
	IsEncrypted int
	ExpirationTime int
	IsReplied int
	IsRepliedAll int
	IsForwarded int
	AddressID string
	Body string
	Header string
	ReplyTo *Email
	Attachments []interface{} // TODO
	Starred int
	Location int
	LabelIDs []string
}

type MessageUpdate struct {
	Message *Message
	IsRead bool
}
