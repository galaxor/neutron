package backend

// Sends messages to email addresses.
type SendBackend interface {
	// Send a message to an e-mail address.
	SendMessagePackage(user string, msg *OutgoingMessage) error
}

// A message that is going to be sent.
// Message.Body MUST be ignored, MessagePackage.Body MUST be used instead.
// The recipient is specified in MessagePackage.Address.
type OutgoingMessage struct {
	*Message
	*MessagePackage

	InReplyTo string
	References string
}
