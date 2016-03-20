package imap

import (
	"errors"
	"net/mail"
	"bytes"
	"strconv"
	"time"
	"mime"
	"mime/multipart"
	"strings"
	"io"
	"io/ioutil"

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

func parseMessageInfo(msg *backend.Message, msgInfo *imap.MessageInfo) {
	msg.ID = strconv.Itoa(int(msgInfo.UID))
	msg.Order = int(msgInfo.Seq)
	msg.Size = int(msgInfo.Size)

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
}

func parseEnvelopeAddress(addr []imap.Field) *backend.Email {
	return &backend.Email{
		Name: imap.AsString(addr[0]),
		Address: imap.AsString(addr[2]) + "@" + imap.AsString(addr[3]),
	}
}

func parseEnvelopeAddressList(list []imap.Field) []*backend.Email {
	emails := make([]*backend.Email, len(list))
	for i, field := range list {
		addr := imap.AsList(field)
		emails[i] = parseEnvelopeAddress(addr)
	}
	return emails
}

func decodeRFC2047Word(word string) string {
	// TODO: mime.WordDecoder cannot handle multiple encoded-words
	// See https://github.com/golang/go/issues/4687#issuecomment-66073826

	dec := new(mime.WordDecoder) // TODO: do not create one decoder per word
	decoded, err := dec.Decode(word)
	if err == nil {
		return decoded
	}
	return word
}

func parseEnvelope(msg *backend.Message, envelope []imap.Field) {
	t, err := time.Parse(time.RFC1123Z, imap.AsString(envelope[0]))
	if err == nil {
		msg.Time = t.Unix()
	}

	msg.Subject = decodeRFC2047Word(imap.AsString(envelope[1]))

	// envelope[2] is From

	senders := imap.AsList(envelope[3])
	if len(senders) > 0 {
		msg.Sender = parseEnvelopeAddress(imap.AsList(senders[0]))
	}

	replyTo := imap.AsList(envelope[4])
	if len(replyTo) > 0 {
		msg.ReplyTo = parseEnvelopeAddress(imap.AsList(replyTo[0]))
	}

	to := imap.AsList(envelope[5])
	msg.ToList = parseEnvelopeAddressList(to)

	cc := imap.AsList(envelope[6])
	msg.CCList = parseEnvelopeAddressList(cc)

	bcc := imap.AsList(envelope[6])
	msg.BCCList = parseEnvelopeAddressList(bcc)

	// envelope[7] is In-Reply-To
	// envelope[8] is Message-Id
}

func parseMessageHeader(msg *backend.Message, header *mail.Header) {
	msg.Subject = decodeRFC2047Word(header.Get("Subject"))

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

	/*body, err := ioutil.ReadAll(m.Body)
	if err == nil && len(body) > 0 {
		msg.Body = string(body)
	}*/
}

func parseMessageBody(msg *backend.Message, m *mail.Message) error {
	mediaType, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	gotType := ""
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(m.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				return err
			}

			partType := p.Header.Get("Content-Type")
			if (strings.HasPrefix(partType, "text/plain") && gotType == "") || strings.HasPrefix(partType, "text/html") {
				gotType = partType
				msg.Body = string(slurp)
			}
		}
	} else {
		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			return err
		}
		msg.Body = string(body)
	}

	return nil
}

func (b *MessagesBackend) GetMessage(user, id string) (msg *backend.Message, err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	set, _ := imap.NewSeqSet(id)
	cmd, err := imap.Wait(c.UIDFetch(set, "UID", "FLAGS", "RFC822.SIZE", "RFC822.HEADER", "RFC822.TEXT"))
	if err != nil {
		return
	}

	rsp := cmd.Data[0]
	msgInfo := rsp.MessageInfo()
	header := imap.AsBytes(msgInfo.Attrs["RFC822.HEADER"])
	body := imap.AsBytes(msgInfo.Attrs["RFC822.TEXT"])
	m, err := mail.ReadMessage(bytes.NewReader(header))
	if err != nil {
		return
	}

	m.Body = bytes.NewReader(body)

	msg = &backend.Message{}
	msg.LabelIDs = []string{getLabelID(c.Mailbox.Name)} // TODO
	parseMessageInfo(msg, msgInfo)
	parseMessageHeader(msg, &m.Header)
	parseMessageBody(msg, m)
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

	cmd, _ := c.Fetch(set, "UID", "FLAGS", "RFC822.SIZE", "ENVELOPE")
	for cmd.InProgress() {
		c.Recv(-1)

		// Process command data
		for _, rsp := range cmd.Data {
			msgInfo := rsp.MessageInfo()
			envelope := imap.AsList(msgInfo.Attrs["ENVELOPE"])

			msg := &backend.Message{}
			msg.LabelIDs = []string{getLabelID(c.Mailbox.Name)} // TODO
			parseMessageInfo(msg, msgInfo)
			parseEnvelope(msg, envelope)

			msgs = append(msgs, msg)
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
