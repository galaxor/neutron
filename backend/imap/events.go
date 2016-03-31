package imap

import (
	"errors"

	"github.com/emersion/neutron/backend"
	"github.com/mxk/go-imap/imap"
)

type Events struct {
	backend.EventsBackend
	conns *conns
	msgs *Messages
}

func (b *Events) DeleteAllEvents(user string) error {
	err := b.EventsBackend.DeleteAllEvents(user)
	if err != nil {
		return err
	}

	return b.conns.disconnect(user)
}

func (b *Events) processUpdate(u *update) error {
	// TODO: support other updates too (EXPUNGE)
	if u.name != "EXISTS" {
		return nil
	}

	user := u.user

	c, unlock, err := b.conns.getConn(user)
	if err != nil {
		return err
	}

	mailbox := c.Mailbox.Name
	seqset, _ := imap.NewSeqSet("")
	seqset.AddNum(u.seqnbr)

	cmd, _, err := wait(c.Fetch(seqset, "UID"))
	unlock()
	if err != nil {
		return err
	}

	if len(cmd.Data) != 1 {
		return errors.New("No such message")
	}

	rsp := cmd.Data[0]
	msgInfo := rsp.MessageInfo()
	uid := msgInfo.UID

	msgId := formatMessageId(mailbox, uid)
	msg, err := b.msgs.GetMessage(user, msgId)
	if err != nil {
		return err
	}

	event := backend.NewMessageDeltaEvent(msg.ID, backend.EventCreate, msg)
	return b.InsertEvent(user, event)
}

func (b *Events) listenUpdates() {
	for {
		u := <-b.conns.updates
		go b.processUpdate(u)
	}
}

func newEvents(conns *conns, events backend.EventsBackend, msgs *Messages) *Events {
	evts := &Events{
		EventsBackend: events,
		conns: conns,
		msgs: msgs,
	}

	go evts.listenUpdates()

	return evts
}
