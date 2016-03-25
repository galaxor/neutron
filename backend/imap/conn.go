package imap

import (
	"sync"
	"errors"
	"strconv"

	"github.com/mxk/go-imap/imap"
)

func (c *Config) Host() string {
	host := c.Hostname
	if c.Port > 0 {
		host += ":" + strconv.Itoa(c.Port)
	}
	return host
}

type conns struct {
	config *Config
	clients map[string]*imap.Client
	locks map[string]sync.Locker
}

func (b *conns) insertConn(user string, conn *imap.Client) {
	b.clients[user] = conn
	b.locks[user] = &sync.Mutex{}
}

func (b *conns) getConn(user string) (*imap.Client, func(), error) {
	lock, ok := b.locks[user]
	if !ok {
		return nil, nil, errors.New("No such user")
	}

	lock.Lock()

	conn, ok := b.clients[user]
	if !ok {
		lock.Unlock()
		return nil, nil, errors.New("No such user")
	}

	state := conn.State()
	if state == imap.Logout || state == imap.Closed {
		delete(b.clients, user)
		delete(b.locks, user)
		lock.Unlock()
		return nil, nil, errors.New("Connection to IMAP server closed")
	}

	return conn, lock.Unlock, nil
}

func newConns(config *Config) *conns {
	return &conns{
		config: config,

		clients: map[string]*imap.Client{},
		locks: map[string]sync.Locker{},
	}
}
