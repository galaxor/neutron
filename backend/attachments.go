package backend

// Stores attachments.
type AttachmentsBackend interface {
	// Get an attachment content.
	ReadAttachment(user, id string) (*Attachment, []byte, error)
	// Insert a new attachment.
	InsertAttachment(user string, attachment *Attachment, contents []byte) (*Attachment, error)
	// Delete an attachment.
	DeleteAttachment(user, id string) error
}

// An attachment.
type Attachment struct {
	ID string
	MessageID string `json:",omitempty"`
	Name string
	Size int
	MIMEType string
	KeyPackets string `json:",omitempty"`
	DataPacket string `json:",omitempty"`
}
