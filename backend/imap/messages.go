package imap

import (
	"errors"
	"encoding/base64"
	"net/mail"
	"bytes"
	"strings"
	"strconv"
	"time"

	"github.com/mxk/go-imap/imap"
	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/util/textproto"
)

type MessagesBackend struct {
	*connBackend

	mailboxes map[string][]*imap.MailboxInfo
}

func formatMessageId(mailbox string, uid uint32) string {
	raw := mailbox + "/" + strconv.Itoa(int(uid))
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func parseMessageId(msgId string) (mailbox string, uid uint32, err error) {
	decoded, err := base64.URLEncoding.DecodeString(msgId)
	if err != nil {
		return
	}

	parts := strings.SplitN(string(decoded), "/", 2)
	if len(parts) != 2 {
		err = errors.New("Invalid message ID: does not contain separator")
		return
	}

	uidInt, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	mailbox = parts[0]
	uid = uint32(uidInt)
	return
}

func parseMessageInfo(msg *backend.Message, msgInfo *imap.MessageInfo) {
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
		msg.LabelIDs = append(msg.LabelIDs, backend.StarredLabel)
	}
	if msgInfo.Flags["\\Draft"] {
		msg.Type = backend.DraftType
	}
}

func parseEnvelopeAddress(addr []imap.Field) *backend.Email {
	return &backend.Email{
		Name: textproto.DecodeWord(imap.AsString(addr[0])),
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

func parseEnvelope(msg *backend.Message, envelope []imap.Field) {
	// TODO: support more formats (see RFC)
	t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 (MST)", imap.AsString(envelope[0]))
	if err != nil {
		t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", imap.AsString(envelope[0]))
	}
	if err == nil {
		msg.Time = t.Unix()
	}

	msg.Subject = textproto.DecodeWord(imap.AsString(envelope[1]))

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

func (b *MessagesBackend) getMailboxes(user string) ([]*imap.MailboxInfo, error) {
	// Mailboxes list already retrieved
	if len(b.mailboxes[user]) > 0 {
		return b.mailboxes[user], nil
	}

	c, unlock, err := b.getConn(user)
	if err != nil {
		return nil, err
	}
	defer unlock()

	// Since the connection was locked, the mailboxes list could now have been
	// retrieved
	if len(b.mailboxes[user]) > 0 {
		return b.mailboxes[user], nil
	}

	cmd, _, err := wait(c.List("", "%"))
	if err != nil {
		return nil, err
	}

	// Retrieve mailboxes info and subscribe to them
	b.mailboxes[user] = make([]*imap.MailboxInfo, len(cmd.Data))
	for i, rsp := range cmd.Data {
		mailboxInfo := rsp.MailboxInfo()
		b.mailboxes[user][i] = mailboxInfo

		_, _, err := wait(c.Subscribe(mailboxInfo.Name))
		if err != nil {
			return nil, err
		}
	}

	return b.mailboxes[user], nil
}

func (b *MessagesBackend) getLabelMailbox(user, label string) (mailbox string, err error) {
	mailboxes, err := b.getMailboxes(user)
	if err != nil {
		return
	}

	mailbox = label
	for _, m := range mailboxes {
		if getLabelID(m.Name) == label {
			mailbox = m.Name
			break
		}
	}

	return
}

func (b *MessagesBackend) selectMailbox(user, mailbox string) (err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	if c.Mailbox == nil || c.Mailbox.Name != mailbox {
		_, err = c.Select(mailbox, false)
		if err != nil {
			return
		}
	}

	return
}

func (b *MessagesBackend) selectLabelMailbox(user, label string) (err error) {
	mailbox, err := b.getLabelMailbox(user, label)
	if err != nil {
		return
	}

	return b.selectMailbox(user, mailbox)
}

func (b *MessagesBackend) GetMessage(user, id string) (msg *backend.Message, err error) {
	mailbox, uid, err := parseMessageId(id)
	if err != nil {
		return
	}

	err = b.selectMailbox(user, mailbox)
	if err != nil {
		return
	}

	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	set, _ := imap.NewSeqSet("")
	set.AddNum(uid)
	cmd, _, err := wait(c.UIDFetch(set, "UID", "FLAGS", "RFC822.SIZE", "RFC822.HEADER", "RFC822.TEXT"))
	if err != nil {
		return
	}

	if len(cmd.Data) != 1 {
		err = errors.New("No such message")
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
	msg.ID = formatMessageId(c.Mailbox.Name, msgInfo.UID)
	msg.LabelIDs = []string{getLabelID(c.Mailbox.Name)}
	msg.Header = string(header)
	parseMessageInfo(msg, msgInfo)
	textproto.ParseMessageHeader(msg, &m.Header)
	textproto.ParseMessageBody(msg, m)
	return
}

func reverseMessagesList(msgs []*backend.Message) {
	n := len(msgs)
	for i := 0; i < n/2; i++ {
		msgs[i], msgs[n-i-1] = msgs[n-i-1], msgs[i]
	}
}

func (b *MessagesBackend) ListMessages(user string, filter *backend.MessagesFilter) (msgs []*backend.Message, total int, err error) {
	if filter.Label == "" {
		err = errors.New("Cannot list messages without specifying a label")
		return
	}

	err = b.selectLabelMailbox(user, filter.Label)
	if err != nil {
		return
	}

	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

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

	cmd, err := c.Fetch(set, "UID", "FLAGS", "RFC822.SIZE", "ENVELOPE")
	if err != nil {
		return
	}

	for cmd.InProgress() {
		c.Recv(-1)

		// Process command data
		for _, rsp := range cmd.Data {
			msgInfo := rsp.MessageInfo()
			envelope := imap.AsList(msgInfo.Attrs["ENVELOPE"])

			msg := &backend.Message{}
			msg.ID = formatMessageId(c.Mailbox.Name, msgInfo.UID)
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

	reverseMessagesList(msgs)
	return
}

func (b *MessagesBackend) CountMessages(user string) (counts []*backend.MessagesCount, err error) {
	mailboxes, err := b.getMailboxes(user)
	if err != nil {
		return
	}

	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	for _, mailbox := range mailboxes {
		cmd, _ := imap.Wait(c.Status(mailbox.Name, "MESSAGES", "UNSEEN"))
		if _, err = cmd.Result(imap.OK); err != nil {
			return
		}

		mailboxStatus := cmd.Data[0].MailboxStatus()

		counts = append(counts, &backend.MessagesCount{
			LabelID: getLabelID(mailboxStatus.Name),
			Total: int(mailboxStatus.Messages),
			Unread: int(mailboxStatus.Unseen),
		})
	}

	return
}

func (b *MessagesBackend) InsertMessage(user string, msg *backend.Message) (inserted *backend.Message, err error) {
	mailbox, err := b.getLabelMailbox(user, backend.DraftLabel)
	if err != nil {
		return
	}

	flags := imap.NewFlagSet("\\Seen", "\\Draft")
	mail := textproto.FormatMessage(msg)
	literal := imap.NewLiteral([]byte(mail))

	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	_, res, err := wait(c.Append(mailbox, flags, nil, literal))
	if err != nil {
		return
	}

	if imap.AsString(res.Fields[0]) != "APPENDUID" {
		err = errors.New("APPEND didn't returned an UID (this is not supported for now)")
		return
	}

	inserted = msg
	inserted.ID = formatMessageId(mailbox, imap.AsNumber(res.Fields[2]))
	return
}

func (b *MessagesBackend) updateMessageFlags(user string, seqset *imap.SeqSet, flag string, value bool) error {
	item := "+FLAGS"
	if !value {
		item = "-FLAGS"
	}

	fields := imap.Field([]imap.Field{imap.Field(flag)})

	c, unlock, err := b.getConn(user)
	if err != nil {
		return err
	}
	defer unlock()

	_, _, err = wait(c.UIDStore(seqset, item, fields))
	if err != nil {
		return err
	}

	return nil
}

// TODO: only supports moving one single message
func (b *MessagesBackend) moveMessages(user string, seqset *imap.SeqSet, mbox string) (uid uint32, err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	_, res, err := wait(c.UIDCopy(seqset, mbox))
	if err != nil {
		return
	}

	if imap.AsString(res.Fields[0]) != "COPYUID" {
		err = errors.New("COPY didn't returned an UID (this is not supported for now)")
		return
	}
	uid = imap.AsNumber(res.Fields[2])

	fields := imap.Field([]imap.Field{imap.Field("\\Deleted")})
	_, _, err = wait(c.UIDStore(seqset, "+FLAGS", fields))
	if err != nil {
		return
	}

	_, _, err = wait(c.Expunge(seqset))
	if err != nil {
		return
	}

	return
}

func (b *MessagesBackend) UpdateMessage(user string, update *backend.MessageUpdate) (msg *backend.Message, err error) {
	mailbox, uid, err := parseMessageId(update.Message.ID)
	if err != nil {
		return
	}

	err = b.selectMailbox(user, mailbox)
	if err != nil {
		return
	}

	seqset, _ := imap.NewSeqSet("")
	seqset.AddNum(uid)

	msg, err = b.GetMessage(user, update.Message.ID)
	if err != nil {
		return
	}

	if update.IsRead {
		err = b.updateMessageFlags(user, seqset, "\\Seen", (update.Message.IsRead == 1))
		if err != nil {
			return
		}
	}

	if update.Starred {
		err = b.updateMessageFlags(user, seqset, "\\Flagged", (update.Message.Starred == 1))
		if err != nil {
			return
		}
	}

	if update.Type {
		err = b.updateMessageFlags(user, seqset, "\\Draft", (update.Message.Type == backend.DraftType))
		if err != nil {
			return
		}
	}

	// TODO: support more scenarios
	var newId string
	if update.LabelIDs == backend.ReplaceLabels && len(update.Message.LabelIDs) == 1 {
		label := update.Message.LabelIDs[0]

		var newMailbox string
		newMailbox, err = b.getLabelMailbox(user, label)
		if err != nil {
			return
		}

		var newUid uint32
		newUid, err = b.moveMessages(user, seqset, newMailbox)
		if err != nil {
			return
		}

		newId = formatMessageId(newMailbox, newUid)
	}

	update.Apply(msg)

	if newId != "" {
		msg.ID = newId
	}

	return
}

func (b *MessagesBackend) DeleteMessage(user, id string) error {
	return errors.New("Not yet implemented")
}

func newMessagesBackend(conn *connBackend) backend.MessagesBackend {
	return &MessagesBackend{
		connBackend: conn,
		mailboxes: map[string][]*imap.MailboxInfo{},
	}
}
