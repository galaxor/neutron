package smtp

import (
	"net/mail"
	"net/textproto"
	"time"

	"github.com/emersion/neutron/backend"
)

func formatEmail(email *backend.Email) string {
	addr := &mail.Address{Name: email.Name, Address: email.Address}
	return addr.String()
}

func getMailHeader(msg *backend.OutgoingMessage) textproto.MIMEHeader {
	h := textproto.MIMEHeader{}

	h.Set("MIME-Version", "1")

	h.Set("Subject", msg.Subject)
	h.Set("From", formatEmail(msg.Sender))
	h.Set("Date", time.Unix(msg.Time, 0).Format(time.RFC1123Z))

	for _, to := range msg.ToList {
		h.Add("To", formatEmail(to))
	}
	for _, cc := range msg.CCList {
		h.Add("Cc", formatEmail(cc))
	}

	if msg.ReplyTo != nil {
		h.Set("Reply-To", formatEmail(msg.ReplyTo))
	}

	if msg.InReplyTo != "" {
		h.Set("In-Reply-To", msg.InReplyTo)
	}
	if msg.References != "" {
		h.Set("References", msg.References)
	}

	return h
}

func fomatHeader(h textproto.MIMEHeader) string {
	output := ""
	for key, values := range h {
		for _, value := range values {
			output += key + ": " + value + "\r\n"
		}
	}
	return output
}
