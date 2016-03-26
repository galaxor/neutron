package imap

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/emersion/neutron/backend"
	"github.com/mxk/go-imap/imap"
)

func (b *Messages) listAttachments(user string, msg *backend.OutgoingMessage) error {
	for _, att := range msg.Attachments {
		att, d, err := b.ReadAttachment(user, att.ID)
		if err != nil {
			return err
		}

		msg.Attachments = append(msg.Attachments, &backend.OutgoingAttachment{
			Attachment: att,
			Data: d,
		})
	}

	return nil
}

func (b *Messages) ReadAttachment(user, id string) (att *backend.Attachment, out []byte, err error) {
	mailbox, uid, partId, err := parseAttachmentId(id)
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

	cmd, _, err := wait(c.UIDFetch(seqset, "BODYSTRUCTURE", "BODY.PEEK["+partId+"]"))
	if err != nil {
		return
	}

	rsp := cmd.Data[0]
	msgInfo := rsp.MessageInfo()
	structure := parseBodyStructure(imap.AsList(msgInfo.Attrs["BODYSTRUCTURE"]))
	body := imap.AsBytes(msgInfo.Attrs["BODY["+partId+"]"])

	part := structure.Get(partId)
	att = part.Attachment()
	r := part.DecodeContent(bytes.NewReader(body))
	out, err = ioutil.ReadAll(r)
	return
}

func (b *Messages) InsertAttachment(user string, attachment *backend.Attachment, data []byte) (*backend.Attachment, error) {
	msgId := attachment.MessageID
	mailbox, uid, err := parseMessageId(msgId)
	if err != nil {
		return nil, err
	}

	msg, err := b.GetMessage(user, msgId)
	if err != nil {
		return nil, err
	}

	outgoing := &backend.OutgoingMessage{
		Message: msg,
	}

	err = b.listAttachments(user, outgoing)
	if err != nil {
		return nil, err
	}

	outgoing.Attachments = append(outgoing.Attachments, &backend.OutgoingAttachment{
		Attachment: attachment,
		Data: data,
	})

	msg, err = b.insertMessage(user, outgoing)
	if err != nil {
		return nil, err
	}

	mailbox, uid, _ = parseMessageId(msg.ID)
	id := strconv.Itoa(len(outgoing.Attachments) + 1 /* body part */)
	attachment.ID = formatAttachmentId(mailbox, uid, id)
	attachment.Size = len(data)
	attachment.MessageID = msg.ID
	return attachment, nil
}

func (b *Messages) DeleteAttachment(user, id string) error {
	return errors.New("Not yet implemented")
}
