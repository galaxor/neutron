package backend

// Stores attachments.
type AttachmentsBackend interface {
	// Get an attachment content.
	ReadAttachment(user, id string) (*Attachment, []byte, error)
}

// An attachment.
type Attachment struct {
	ID string
	Name string
	Size int
	MIMEType string
	KeyPackets string `json:",omitempty"`
	DataPacket string `json:",omitempty"`
}
