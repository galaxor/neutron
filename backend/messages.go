package backend

// Stores messages data.
type MessagesBackend interface {
	// Get a message.
	GetMessage(user, id string) (*Message, error)
	// List all user's messages. A message filter can be provided.
	ListMessages(user string, filter *MessagesFilter) ([]*Message, int, error)
	// Count all user's messages by label.
	CountMessages(user string) ([]*MessagesCount, error)
	// Insert a new message.
	InsertMessage(user string, msg *Message) (*Message, error)
	// Update an existing message.
	UpdateMessage(user string, update *MessageUpdate) (*Message, error)
	// Permanently delete a message.
	DeleteMessage(user, id string) error
}

// Sends messages to email addresses.
type SendBackend interface {
	// Send a message to an e-mail address.
	SendMessagePackage(user string, msg *Message, pkg *MessagePackage) error
}

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

// Apply this update on a message.
func (update *MessageUpdate) Apply(msg *Message) {
	updated := update.Message

	if updated.ID != msg.ID {
		panic("Cannot apply update on a message with a different ID")
	}

	if update.ToList {
		msg.ToList = updated.ToList
	}
	if update.CCList {
		msg.CCList = updated.CCList
	}
	if update.BCCList {
		msg.BCCList = updated.BCCList
	}
	if update.Subject {
		msg.Subject = updated.Subject
	}
	if update.IsRead {
		msg.IsRead = updated.IsRead
	}
	if update.Type {
		msg.Type = updated.Type
	}
	if update.AddressID {
		msg.AddressID = updated.AddressID
	}
	if update.Body {
		msg.Body = updated.Body
	}
	if update.Time {
		msg.Time = updated.Time
	}

	if update.LabelIDs != KeepLabels {
		switch update.LabelIDs {
		case ReplaceLabels:
			msg.LabelIDs = updated.LabelIDs
		case AddLabels:
			for _, lblToAdd := range updated.LabelIDs {
				found := false
				for _, lbl := range msg.LabelIDs {
					if lbl == lblToAdd {
						found = true
						break
					}
				}
				if !found {
					msg.LabelIDs = append(msg.LabelIDs, lblToAdd)
				}
			}
		case RemoveLabels:
			labels := []string{}
			for _, lbl := range msg.LabelIDs {
				found := false
				for _, lblToRemove := range updated.LabelIDs {
					if lbl == lblToRemove {
						found = true
						break
					}
				}
				if !found {
					labels = append(labels, lbl)
				}
			}
			msg.LabelIDs = labels
		}
	}
}
