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
	HasAttachment int
	NumAttachments int
	IsEncrypted int
	ExpirationTime int
	IsReplied int
	IsRepliedAll int
	IsForwarded int
	AddressID string
	Body string `json:",omitempty"`
	Header string `json:",omitempty"`
	ReplyTo *Email
	Attachments []*Attachment
	Starred int
	Location int
	LabelIDs []string
}

// An email contains an address and a name.
type Email struct {
	Name string
	Address string
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

// Messages locations.
const (
	InboxLocation int = 0
	DraftLocation = 1
	SentLocation = 2
	TrashLocation = 3
	SpamLocation = 4
	ArchiveLocation = 6
)

type MessagePackage struct {
	Address string
	Type int
	Body string
	KeyPackets []string
}

// Contains message counts for one label.
type MessagesCount struct {
	LabelID string
	Total int
	Unread int
}

// Contains a summary of messages counts per location and label.
type MessagesTotal struct {
	Locations []*LocationTotal
	Labels []*LabelTotal
	Starred int
}

type LocationTotal struct {
	Location int
	Count int
}

type LabelTotal struct {
	LabelID string
	Count int
}

func addCountToTotal(totals *MessagesTotal, label string, count int) {
	if count == 0 {
		return
	}

	if label == StarredLabel {
		totals.Starred += count
		return
	}

	location := -1
	switch label {
	case InboxLabel:
		location = InboxLocation
	case DraftLabel:
		location = DraftLocation
	case SentLabel:
		location = SentLocation
	case TrashLabel:
		location = TrashLocation
	case SpamLabel:
		location = SpamLocation
	case ArchiveLabel:
		location = ArchiveLocation
	}

	if location != -1 { // A system label that has a corresponding location
		found := false
		for _, t := range totals.Locations {
			if t.Location == location {
				found = true
				t.Count += count
				break
			}
		}

		if !found {
			totals.Locations = append(totals.Locations, &LocationTotal{
				Location: location,
				Count: count,
			})
		}
	} else { // Just a regular label
		found := false
		for _, t := range totals.Labels {
			if t.LabelID == label {
				found = true
				t.Count += count
				break
			}
		}

		if !found {
			totals.Labels = append(totals.Labels, &LabelTotal{
				LabelID: label,
				Count: count,
			})
		}
	}
}

func MessagesTotalFromCounts(counts []*MessagesCount) (totals *MessagesTotal, unread *MessagesTotal) {
	totals = &MessagesTotal{}
	unread = &MessagesTotal{}

	for _, count := range counts {
		addCountToTotal(totals, count.LabelID, count.Total)
		addCountToTotal(unread, count.LabelID, count.Unread)
	}

	return
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
	Begin int64 // Timestamp
	End int64 // Timestamp
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
