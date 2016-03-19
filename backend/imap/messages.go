package imap

import (
	"errors"
	"net/mail"
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/mxk/go-imap/imap"
	"github.com/emersion/neutron/backend"
)

type MessagesBackend struct {
	*connBackend
}

func populateMessage(msg *backend.Message) {
	if msg.ToList == nil {
		msg.ToList = []*backend.Email{}
	}
	if msg.CCList == nil {
		msg.CCList = []*backend.Email{}
	}
	if msg.BCCList == nil {
		msg.BCCList = []*backend.Email{}
	}
	if msg.Attachments == nil {
		msg.Attachments = []*backend.Attachment{}
	}
	if msg.LabelIDs == nil {
		msg.LabelIDs = []string{}
	}

	if msg.Sender != nil {
		msg.SenderAddress = msg.Sender.Address
		msg.SenderName = msg.Sender.Name
	}

	if backend.IsEncrypted(msg.Body) {
		msg.IsEncrypted = 1
	}
}

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

func getEmail(addr *mail.Address) *backend.Email {
	return &backend.Email{
		Name: addr.Name,
		Address: addr.Address,
	}
}

func getMessage(msgInfo *imap.MessageInfo, b []byte) *backend.Message {
	m, err := mail.ReadMessage(bytes.NewReader(b))
	if m == nil || err != nil {
		return nil
	}

	header := m.Header

	msg := &backend.Message{
		ID: strconv.Itoa(int(msgInfo.UID)),
		Order: int(msgInfo.Seq),
		Subject: header.Get("Subject"),
		Size: int(msgInfo.Size),
		LabelIDs: []string{backend.InboxLabel}, // TODO
	}

	if msgInfo.Flags["\\Seen"] {
		msg.IsRead = 1
	}
	if msgInfo.Flags["\\Answered"] {
		msg.IsReplied = 1
	}
	if msgInfo.Flags["\\Flagged"] {
		msg.Starred = 1
	}
	if msgInfo.Flags["\\Draft"] {
		msg.Type = backend.DraftType
	}

	from, err := header.AddressList("From")
	if err == nil && len(from) > 0 {
		msg.Sender = getEmail(from[0])
	}

	to, err := header.AddressList("To")
	if err == nil {
		for _, addr := range to {
			msg.ToList = append(msg.ToList, getEmail(addr))
		}
	}

	// TODO: CCList, BCCList

	replyTo, err := header.AddressList("From")
	if err == nil && len(replyTo) > 0 {
		msg.ReplyTo = getEmail(replyTo[0])
	}

	time, err := header.Date()
	if err == nil {
		msg.Time = time.Unix()
	}

	body, err := ioutil.ReadAll(m.Body)
	if err == nil {
		msg.Body = string(body)
	}

	return msg
}

func (b *MessagesBackend) GetMessage(user, id string) (msg *backend.Message, err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	set, _ := imap.NewSeqSet(id)
	cmd, err := imap.Wait(c.UIDFetch(set, "UID", "FLAGS", "RFC822.SIZE", "RFC822.HEADER", "BODY"))
	if err != nil {
		return
	}

	rsp := cmd.Data[0]
	msgInfo := rsp.MessageInfo()
	header := imap.AsBytes(msgInfo.Attrs["RFC822.HEADER"])
	msg = getMessage(msgInfo, header)
	if msg == nil {
		err = errors.New("Cannot parse message headers")
		return
	}
	msg.Body = imap.AsString(msgInfo.Attrs["BODY"])
	return
}

func (b *MessagesBackend) ListMessages(user string, filter *backend.MessagesFilter) (msgs []*backend.Message, total int, err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	mailbox := "INBOX" // TODO: use filter.Label
	filter.Limit = 10

	c.Select(mailbox, true)

	total = int(c.Mailbox.Messages) // TODO: not filtered

	set, _ := imap.NewSeqSet("")
	if filter.Limit > 0 && filter.Page >= 0 {
		from := filter.Limit * filter.Page
		to := filter.Limit * (filter.Page + 1)

		if uint32(to) < c.Mailbox.Messages {
			set.AddRange(c.Mailbox.Messages - uint32(from), c.Mailbox.Messages - uint32(to))
		} else {
			set.Add("1:*")
		}
	} else {
		set.Add("1:*")
	}

	cmd, _ := c.Fetch(set, "UID", "FLAGS", "RFC822.SIZE", "RFC822.HEADER")
	for cmd.InProgress() {
		c.Recv(-1)

		// Process command data
		for _, rsp := range cmd.Data {
			msgInfo := rsp.MessageInfo()
			header := imap.AsBytes(msgInfo.Attrs["RFC822.HEADER"])

			msg := getMessage(msgInfo, header)
			if msg != nil {
				msgs = append(msgs, msg)
			}
		}
		cmd.Data = nil
	}

	c.Data = nil

	// Check command completion status
	if _, err = cmd.Result(imap.OK); err != nil {
		return
	}

	return
}

func (b *MessagesBackend) CountMessages(user string) (counts []*backend.MessagesCount, err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	cmd, _ := imap.Wait(c.List("", "%"))

	for _, rsp := range cmd.Data {
		mailboxInfo := rsp.MailboxInfo()

		cmd, _ = imap.Wait(c.Status(mailboxInfo.Name, "MESSAGES", "UNSEEN"))
		mailboxStatus := cmd.Data[0].MailboxStatus()

		counts = append(counts, &backend.MessagesCount{
			LabelID: getLabelID(mailboxStatus.Name),
			Total: int(mailboxStatus.Messages),
			Unread: int(mailboxStatus.Unseen),
		})
	}

	return
}

func (b *MessagesBackend) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	return nil, errors.New("Not yet implemented")
}

func (b *MessagesBackend) UpdateMessage(user string, update *backend.MessageUpdate) (*backend.Message, error) {
	return nil, errors.New("Not yet implemented")
}

func (b *MessagesBackend) DeleteMessage(user, id string) error {
	return errors.New("Not yet implemented")
}
