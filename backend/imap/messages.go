package imap

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/memory"
	"github.com/emersion/neutron/backend/util/textproto"
	"github.com/mxk/go-imap/imap"
)

type updatableAttachments interface {
	backend.AttachmentsBackend
	UpdateAttachmentMessage(user, id, msgId string) error
}

type Messages struct {
	*conns
	tmpAtts updatableAttachments
}

func formatAttachmentId(mailbox string, uid uint32, part string) string {
	raw := mailbox + "/" + strconv.Itoa(int(uid))
	if part != "" {
		raw += "#" + part
	}
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func formatMessageId(mailbox string, uid uint32) string {
	return formatAttachmentId(mailbox, uid, "")
}

func parseAttachmentId(id string) (mailbox string, uid uint32, part string, err error) {
	decoded, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return
	}

	fstParts := strings.SplitN(string(decoded), "/", 2)
	if len(fstParts) != 2 {
		err = errors.New("Invalid message ID: does not contain separator")
		return
	}
	sndParts := strings.SplitN(fstParts[1], "#", 2)

	uidInt, err := strconv.Atoi(sndParts[0])
	if err != nil {
		return
	}

	mailbox = fstParts[0]
	uid = uint32(uidInt)

	if len(sndParts) == 2 {
		part = sndParts[1]
	}
	return
}

func parseMessageId(id string) (mailbox string, uid uint32, err error) {
	mailbox, uid, _, err = parseAttachmentId(id)
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
		Name:    textproto.DecodeWord(imap.AsString(addr[0])),
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

func parseBodyStructureParams(params []imap.Field) map[string]string {
	result := map[string]string{}

	for i := 0; i < len(params); i += 2 {
		key := imap.AsString(params[i])
		val := imap.AsString(params[i+1])

		result[key] = val
	}

	return result
}

func parseBodyStructure(structure []imap.Field) *textproto.BodyStructure {
	var parse func(structure []imap.Field, id string) *textproto.BodyStructure
	parse = func(structure []imap.Field, id string) *textproto.BodyStructure {
		if imap.TypeOf(structure[0]) == imap.QuotedString {
			if id == "" {
				id = "1"
			}

			// Not a MIME message
			return &textproto.BodyStructure{
				ID:                 id,
				Type:               imap.AsString(structure[0]),
				SubType:            imap.AsString(structure[1]),
				Params:             parseBodyStructureParams(imap.AsList(structure[2])),
				ContentId:          imap.AsString(structure[3]),
				ContentDescription: imap.AsString(structure[4]),
				ContentEncoding:    imap.AsString(structure[5]),
				Size:               int(imap.AsNumber(structure[6])),
			}
		}

		var processedUntil int
		var children []*textproto.BodyStructure
		for i, field := range structure {
			if imap.TypeOf(field) != imap.List {
				processedUntil = i
				break
			}

			childId := strconv.Itoa(i + 1)
			if id != "" {
				childId = id + "." + childId
			}

			child := parse(imap.AsList(field), childId)
			children = append(children, child)
		}

		return &textproto.BodyStructure{
			ID:       id,
			Type:     "multipart",
			SubType:  imap.AsString(structure[processedUntil]),
			Children: children,
		}
	}

	return parse(structure, "")
}

func (b *Messages) GetMessage(user, id string) (msg *backend.Message, err error) {
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

	seqset, _ := imap.NewSeqSet("")
	seqset.AddNum(uid)

	// Get message metadata

	cmd, _, err := wait(c.UIDFetch(seqset, "FLAGS", "RFC822.SIZE", "RFC822.HEADER", "BODYSTRUCTURE"))
	if err != nil {
		return
	}

	if len(cmd.Data) != 1 {
		err = errors.New("No such message")
		return
	}

	rsp := cmd.Data[0]
	msgInfo := rsp.MessageInfo()
	structure := parseBodyStructure(imap.AsList(msgInfo.Attrs["BODYSTRUCTURE"]))

	header := imap.AsBytes(msgInfo.Attrs["RFC822.HEADER"])
	m, err := mail.ReadMessage(bytes.NewReader(header))
	if err != nil {
		return
	}

	msg = &backend.Message{}
	msg.ID = id
	msg.LabelIDs = []string{getLabelID(c.Mailbox.Name)}
	msg.Header = string(header)
	parseMessageInfo(msg, msgInfo)
	textproto.ParseMessageHeader(msg, &m.Header)
	textproto.ParseMessageStructure(msg, structure)

	// Get message content

	preferred := structure.GetPreferredPart()
	cmd, _, err = wait(c.UIDFetch(seqset, "BODY.PEEK["+preferred.ID+"]"))
	if err != nil {
		return
	}

	rsp = cmd.Data[0]
	msgInfo = rsp.MessageInfo()
	body := imap.AsBytes(msgInfo.Attrs["BODY["+preferred.ID+"]"])

	r := preferred.DecodeContent(bytes.NewReader(body))
	slurp, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	msg.Body = string(slurp)

	// Check if some attachments are encrypted

	for _, att := range msg.Attachments {
		if att.MIMEType == "application/pgp" {
			// TODO: get correct key packet length
			// TODO: attachment is assumed to be base64-encoded
			// TODO: make sure this is a normalized msgInfo.Attrs key
			cmd, _, err = wait(c.UIDFetch(seqset, "BODY.PEEK["+att.ID+"]<0.2048>"))
			if err != nil {
				return
			}

			rsp = cmd.Data[0]
			msgInfo = rsp.MessageInfo()
			encryptedKeyPacket := imap.AsBytes(msgInfo.Attrs["BODY["+att.ID+"]<0>"])
			att.KeyPackets = string(encryptedKeyPacket)
		}

		att.ID = formatAttachmentId(mailbox, uid, att.ID)
	}

	tmpAtts, err := b.tmpAtts.ListAttachments(user, id)
	if err == nil {
		msg.Attachments = append(msg.Attachments, tmpAtts...)
	}

	return
}

func reverseMessagesList(msgs []*backend.Message) {
	n := len(msgs)
	for i := 0; i < n/2; i++ {
		msgs[i], msgs[n-i-1] = msgs[n-i-1], msgs[i]
	}
}

func (b *Messages) ListMessages(user string, filter *backend.MessagesFilter) (msgs []*backend.Message, total int, err error) {
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

	total = int(c.Mailbox.Messages)

	// No messages to fetch
	if c.Mailbox.Messages == 0 {
		return
	}

	set, _ := imap.NewSeqSet("")
	if filter.Limit > 0 && filter.Page >= 0 {
		from := filter.Limit * filter.Page
		to := filter.Limit * (filter.Page + 1)

		if uint32(to) < c.Mailbox.Messages {
			set.AddRange(c.Mailbox.Messages-uint32(from), c.Mailbox.Messages-uint32(to))
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
			msg.LabelIDs = []string{getLabelID(c.Mailbox.Name)}
			parseMessageInfo(msg, msgInfo)
			parseEnvelope(msg, envelope)

			tmpAtts, err := b.tmpAtts.ListAttachments(user, msg.ID)
			if err == nil {
				msg.Attachments = append(msg.Attachments, tmpAtts...)
			}

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

func (b *Messages) CountMessages(user string) (counts []*backend.MessagesCount, err error) {
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
		cmd, _, err := wait(c.Status(mailbox.Name, "MESSAGES", "UNSEEN"))
		if err != nil {
			return nil, err
		}

		mailboxStatus := cmd.Data[0].MailboxStatus()

		counts = append(counts, &backend.MessagesCount{
			LabelID: getLabelID(mailboxStatus.Name),
			Total:   int(mailboxStatus.Messages),
			Unread:  int(mailboxStatus.Unseen),
		})
	}

	return
}

func (b *Messages) insertMessage(user, mailbox string, flags imap.FlagSet, mail []byte) (uid uint32, err error) {
	literal := imap.NewLiteral(mail)

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

	uid = imap.AsNumber(res.Fields[2])
	return
}

func (b *Messages) InsertMessage(user string, msg *backend.Message) (inserted *backend.Message, err error) {
	mailbox, err := b.getLabelMailbox(user, backend.DraftLabel)
	if err != nil {
		return
	}

	flags := imap.NewFlagSet("\\Seen", "\\Draft")
	mail := textproto.FormatMessage(msg)

	uid, err := b.insertMessage(user, mailbox, flags, []byte(mail))

	inserted = msg
	inserted.ID = formatMessageId(mailbox, uid)
	return
}

func (b *Messages) updateMessageFlags(user string, seqset *imap.SeqSet, flag string, value bool) error {
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

func (b *Messages) deleteMessages(user string, seqset *imap.SeqSet) (err error) {
	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

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

// TODO: only supports moving one single message
func (b *Messages) copyMessages(user string, seqset *imap.SeqSet, mbox string) (uid uint32, err error) {
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
	return
}

func (b *Messages) moveMessages(user string, seqset *imap.SeqSet, mbox string) (uid uint32, err error) {
	uid, err = b.copyMessages(user, seqset, mbox)
	if err != nil {
		return
	}

	err = b.deleteMessages(user, seqset)
	return
}

func (b *Messages) UpdateMessage(user string, update *backend.MessageUpdate) (msg *backend.Message, err error) {
	// Retrieve message from mailbox
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

	// Apply update to message
	update.Apply(msg)

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

	// Mark a message as sent
	// Correctly handle temporary attachments
	if update.Type && update.Message.Type == backend.SentType {
		// Retrieve temporary attachments
		tmpAtts, _ := b.tmpAtts.ListAttachments(user, update.Message.ID)

		// Insert temporary attachments to message and send it to the server

		outgoing := &backend.OutgoingMessage{
			Message: msg,
		}

		for _, att := range tmpAtts {
			var d []byte
			att, d, err = b.tmpAtts.ReadAttachment(user, att.ID)

			outgoing.Attachments = append(outgoing.Attachments, &backend.OutgoingAttachment{
				Attachment: att,
				Data: d,
			})
		}

		var mailbox string
		mailbox, err = b.getLabelMailbox(user, backend.SentLabel)
		if err != nil {
			return
		}

		flags := imap.NewFlagSet("\\Seen")
		mail := textproto.FormatOutgoingMessage(outgoing)

		var uid uint32
		uid, err = b.insertMessage(user, mailbox, flags, []byte(mail))
		if err != nil {
			return
		}

		msg.ID = formatMessageId(mailbox, uid)

		// Remove temporary attachments
		for _, att := range tmpAtts {
			b.tmpAtts.DeleteAttachment(user, att.ID)
		}
	} else if update.ToList || update.CCList || update.BCCList || update.Subject || update.AddressID || update.Body || update.Time {
		// If one of those is modified, we have to re-send the whole message to the server

		// The message ID will change
		oldId := msg.ID
		msg.ID = "" // Will be overwritten

		// Insert the updated message
		msg, err = b.InsertMessage(user, msg)
		if err != nil {
			return
		}

		// Delete the old message
		err = b.deleteMessages(user, seqset)
		if err != nil {
			return
		}

		// Update temporary attachments message ID
		tmpAtts, _ := b.tmpAtts.ListAttachments(user, oldId)
		for _, att := range tmpAtts {
			b.tmpAtts.UpdateAttachmentMessage(user, att.ID, msg.ID)
		}
	} else if update.LabelIDs != backend.RemoveLabels && len(update.Message.LabelIDs) == 1 {
		// Move the message from its mailbox to another one
		// TODO: support more scenarios

		msg.LabelIDs = update.Message.LabelIDs
		label := msg.LabelIDs[0]

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

		// Update message ID
		oldId := msg.ID
		msg.ID = formatMessageId(newMailbox, newUid)

		// Update temporary attachments message ID
		tmpAtts, _ := b.tmpAtts.ListAttachments(user, oldId)
		for _, att := range tmpAtts {
			b.tmpAtts.UpdateAttachmentMessage(user, att.ID, msg.ID)
		}
	}

	return
}

func (b *Messages) DeleteMessage(user, id string) (err error) {
	mailbox, uid, err := parseMessageId(id)
	if err != nil {
		return
	}

	err = b.selectMailbox(user, mailbox)
	if err != nil {
		return
	}

	seqset, _ := imap.NewSeqSet("")
	seqset.AddNum(uid)

	err = b.deleteMessages(user, seqset)
	return
}

func newMessages(conns *conns) *Messages {
	return &Messages{
		conns:     conns,
		tmpAtts:   memory.NewAttachments().(updatableAttachments),
	}
}
