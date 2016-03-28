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
	mailboxes map[string][]*imap.MailboxInfo
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

func (b *conns) getMailboxes(user string) ([]*imap.MailboxInfo, error) {
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

func (b *conns) getLabelMailbox(user, label string) (mailbox string, err error) {
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

func (b *conns) selectMailbox(user, mailbox string) (err error) {
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

func (b *conns) selectLabelMailbox(user, label string) (err error) {
	mailbox, err := b.getLabelMailbox(user, label)
	if err != nil {
		return
	}

	return b.selectMailbox(user, mailbox)
}

func newConns(config *Config) *conns {
	return &conns{
		config: config,

		clients: map[string]*imap.Client{},
		passwords: map[string]string{},
		locks: map[string]sync.Locker{},
		mailboxes: map[string][]*imap.MailboxInfo{},
	}
}
