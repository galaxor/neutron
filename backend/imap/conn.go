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
	passwords map[string]string
	locks map[string]sync.Locker
}

func (b *conns) connect(username, password string) (email string, err error) {
	c, err := imap.DialTLS(b.config.Host(), nil)
	if err != nil {
		return
	}

	email = username + b.config.Suffix
	_, err = c.Login(email, password)
	if err != nil {
		return
	}

	b.passwords[username] = password
	b.clients[username] = c
	b.locks[username] = &sync.Mutex{}
	return
}

func (b *conns) getConn(user string) (*imap.Client, func(), error) {
	lock, ok := b.locks[user]
	if !ok {
		return nil, nil, errors.New("No such user")
	}

	lock.Lock()

	c, ok := b.clients[user]
	if !ok {
		lock.Unlock()
		return nil, nil, errors.New("No such user")
	}

	state := c.State()
	if state == imap.Logout || state == imap.Closed {
		// Connection closed, reconnect
		_, err := b.connect(user, b.passwords[user])
		if err != nil {
			delete(b.clients, user)
			delete(b.passwords, user)
			delete(b.locks, user)
			lock.Unlock()
			return nil, nil, err
		}

		c = b.clients[user]
	}

	return c, lock.Unlock, nil
}

// Allow other backends (e.g. a SMTP backend) to access users' password.
func (b *conns) GetPassword(user string) (string, error) {
	if password, ok := b.passwords[user]; ok {
		return password, nil
	}
	return "", errors.New("No password stored for such user")
}

func newConns(config *Config) *conns {
	return &conns{
		config: config,

		clients: map[string]*imap.Client{},
		passwords: map[string]string{},
		locks: map[string]sync.Locker{},
	}
}
