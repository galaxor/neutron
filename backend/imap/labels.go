package imap

import (
	"github.com/emersion/neutron/backend"
)

func getLabelID(mailbox string) string {
	lbl := mailbox
	switch lbl {
	case "INBOX":
		lbl = backend.InboxLabel
	case "Draft", "Drafts":
		lbl = backend.DraftsLabel
	case "Sent":
		lbl = backend.SentLabel
	case "Trash":
		lbl = backend.TrashLabel
	case "Spam", "Junk":
		lbl = backend.SpamLabel
	case "Archive", "Archives":
		lbl = backend.ArchiveLabel
	case "Important", "Starred":
		lbl = backend.StarredLabel
	}
	return lbl
}
