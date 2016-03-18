package backend

// A message.
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
	Time int64
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
	Attachments []*Attachment
	Starred int
	Location int
	LabelIDs []string
}

// Message types.
const (
	DraftType int = 1
	SentType = 2
	SentToMyselfType = 3
)

// Message encryption types.
const (
	Unencrypted int = 0
	EndToEndEncryptedInternal = 1
	EncryptedExternal = 2
	EndToEndEncryptedExternal = 3
	StoredEncryptedExternal = 4
	StoredEncrypted = 5
	EndToEndEncryptedExternalReply = 6
	EncryptedPgp = 7
	EncryptedPgpMime = 8
)

type Attachment struct {} // TODO

type MessagePackage struct {
	Address string
	Type int
	Body string
	KeyPackets []interface{} // TODO
}

// Contains message counts for one label.
type MessagesCount struct {
	LabelID string
	Total int
	Unread int
}

// Contains fields to filter messages.
type MessagesFilter struct {
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

// A request to update a message.
// Fields set to true will be updated with values in Message.
type MessageUpdate struct {
	Message *Message
	ToList bool
	CCList bool
	BCCList bool
	Subject bool
	IsRead bool
	Type bool
	AddressID bool
	Body bool
	Time bool
	Starred bool
	LabelIDs LabelsOperation
}

// The operation to apply to labels.
type LabelsOperation int

const (
	KeepLabels LabelsOperation = iota // Do nothing
	ReplaceLabels // Replace current labels with new ones
	AddLabels // Add new labels to current ones
	RemoveLabels // Remove specified labels from current ones
)
