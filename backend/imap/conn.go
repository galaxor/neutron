package imap

import (
	"sync"
	"errors"
	"strconv"

	"github.com/mxk/go-imap/imap"
)

type Config struct {
	Hostname string
	Port int
	Suffix string
}

func (c *Config) Host() string {
	host := c.Hostname
	if c.Port > 0 {
		host += ":" + strconv.Itoa(c.Port)
	}
	return host
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
		lock.Unlock()
		return nil, nil, errors.New("No such user")
	}

	state := conn.State()
	if state == imap.Logout || state == imap.Closed {
		delete(b.conns, user)
		delete(b.locks, user)
		lock.Unlock()
		return nil, nil, errors.New("Connection to IMAP server closed")
	}

	return conn, lock.Unlock, nil
}

func newConnBackend(config *Config) *connBackend {
	return &connBackend{
		config: config,

		conns: map[string]*imap.Client{},
		locks: map[string]sync.Locker{},
	}
}
