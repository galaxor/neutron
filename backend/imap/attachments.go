package imap

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/emersion/neutron/backend"
	"github.com/mxk/go-imap/imap"
)

func (b *Messages) ListAttachments(user, msg string) ([]*backend.Attachment, error) {
	return nil, errors.New("Not yet implemented")
}

func (b *Messages) ReadAttachment(user, id string) (att *backend.Attachment, out []byte, err error) {
	// First, try to get attachment from temporary backend
	att, out, err = b.tmpAtts.ReadAttachment(user, id)
	if err == nil {
		return
	}

	// Not found in tmp backend, get it from the server

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
	return b.tmpAtts.InsertAttachment(user, attachment, data)
}

func (b *Messages) DeleteAttachment(user, id string) error {
	return b.tmpAtts.DeleteAttachment(user, id)
}
