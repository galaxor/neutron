package imap

import (
	"sync"
	"errors"

	"github.com/mxk/go-imap/imap"
	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/memory"
	"github.com/emersion/neutron/backend/util"
)

type Backend struct {
	backend.DomainsBackend
	backend.ContactsBackend
	backend.LabelsBackend
	backend.ConversationsBackend
	backend.SendBackend
	backend.EventsBackend

	*connBackend

	users map[string]*backend.User
}

func New() backend.Backend {
	bkd := &Backend{
		users: map[string]*backend.User{},
		connBackend: newConnBackend(),
	}

	messages := newMessagesBackend(bkd.connBackend)
	conversations := util.NewDummyConversationsBackend(messages)

	// TODO: do not use memory backends
	bkd.EventsBackend = memory.NewEventsBackend()
	bkd.DomainsBackend = memory.NewDomainsBackend()
	bkd.ContactsBackend = util.NewEventedContactsBackend(memory.NewContactsBackend(), bkd.EventsBackend)
	bkd.LabelsBackend = util.NewEventedLabelsBackend(memory.NewLabelsBackend(), bkd.EventsBackend)
	bkd.ConversationsBackend = util.NewEventedConversationsBackend(conversations, bkd.EventsBackend)
	bkd.SendBackend = util.NewEchoSendBackend(bkd.ConversationsBackend)

	return bkd
}


type Config struct {
	Host string
	Suffix string
}

type connBackend struct {
	config *Config
	conns map[string]*imap.Client
	locks map[string]sync.Locker
}

func (b *connBackend) insertConn(user string, conn *imap.Client) {
	b.conns[user] = conn
	b.locks[user] = &sync.Mutex{}
}

func (b *connBackend) getConn(user string) (*imap.Client, func(), error) {
	lock, ok := b.locks[user]
	if !ok {
		return nil, nil, errors.New("No such user")
	}

	lock.Lock()

	conn, ok := b.conns[user]
	if !ok {
		return nil, nil, errors.New("No such user")
	}

	return conn, lock.Unlock, nil
}

func newConnBackend() *connBackend {
	return &connBackend{
		// TODO: make this configurable
		config: &Config{
			Host: "mail.gandi.net",
			Suffix: "@emersion.fr",
		},

		conns: map[string]*imap.Client{},
		locks: map[string]sync.Locker{},
	}
}
