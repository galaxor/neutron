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
	"log"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"github.com/mxk/go-imap/imap"
	"github.com/emersion/neutron/backend"
)

type MessagesBackend struct {
	*connBackend

	mailboxes map[string][]*imap.MailboxInfo
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

func decodeRFC2047Word(word string) string {
	// TODO: mime.WordDecoder cannot handle multiple encoded-words
	// See https://github.com/golang/go/issues/4687#issuecomment-66073826

	dec := new(mime.WordDecoder) // TODO: do not create one decoder per word
	decoded, err := dec.DecodeHeader(word)
	if err == nil {
		return decoded
	}
	return word
}

func parseEnvelopeAddress(addr []imap.Field) *backend.Email {
	return &backend.Email{
		Name: decodeRFC2047Word(imap.AsString(addr[0])),
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
	t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 (MST)", imap.AsString(envelope[0]))
	if err != nil {
		t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", imap.AsString(envelope[0]))
	}
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

func getEmail(addr *mail.Address) *backend.Email {
	return &backend.Email{
		Name: addr.Name,
		Address: addr.Address,
	}
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

func decodeBytes(b []byte, charset string) []byte {
	var enc encoding.Encoding
	switch strings.ToLower(charset) {
	case "iso-8859-1":
		enc = charmap.ISO8859_1
	case "windows-1252":
		enc = charmap.Windows1252
	case "utf-8":
		// Nothing to do
	default:
		if charset != "" {
			log.Println("WARN: unsupported charset:", charset)
		}
	}
	if enc != nil {
		b, _ = enc.NewDecoder().Bytes(b)
	}
	return b
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

			mediaType, params, err = mime.ParseMediaType(p.Header.Get("Content-Type"))
			if (mediaType == "text/plain" && gotType == "") || mediaType == "text/html" {
				gotType = mediaType
				msg.Body = string(decodeBytes(slurp, params["charset"]))
			}
		}
	} else {
		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			return err
		}
		msg.Body = string(decodeBytes(body, params["charset"]))
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
	msg.LabelIDs = []string{getLabelID(c.Mailbox.Name)} // TODO
	msg.Header = string(header)
	parseMessageInfo(msg, msgInfo)
	parseMessageHeader(msg, &m.Header)
	parseMessageBody(msg, m)
	return
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

	// Since the connection was locked, the mailboxes list could new have been
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

	mailbox, err := b.getLabelMailbox(user, filter.Label)
	if err != nil {
		return
	}

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

func (b *MessagesBackend) InsertMessage(user string, msg *backend.Message) (*backend.Message, error) {
	return nil, errors.New("Not yet implemented")
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

func (b *MessagesBackend) moveMessages(user string, seqset *imap.SeqSet, mbox string) error {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return err
	}
	defer unlock()

	_, _, err = wait(c.UIDCopy(seqset, mbox))
	if err != nil {
		return err
	}

	fields := imap.Field([]imap.Field{imap.Field("\\Deleted")})
	_, _, err = wait(c.UIDStore(seqset, "+FLAGS", fields))
	if err != nil {
		return err
	}

	_, _, err = wait(c.Expunge(seqset))
	if err != nil {
		return err
	}

	return nil
}

func (b *MessagesBackend) UpdateMessage(user string, update *backend.MessageUpdate) (msg *backend.Message, err error) {
	id := update.Message.ID
	seqset, _ := imap.NewSeqSet(id)

	msg, err = b.GetMessage(user, id)
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
	if update.LabelIDs == backend.ReplaceLabels && len(update.Message.LabelIDs) == 1 {
		label := update.Message.LabelIDs[0]

		var mbox string
		mbox, err = b.getLabelMailbox(user, label)
		if err != nil {
			return
		}

		err = b.moveMessages(user, seqset, mbox)
		if err != nil {
			return
		}
	}

	update.Apply(msg)
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
