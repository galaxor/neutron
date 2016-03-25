package imap

import (
	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/util"
)

type Config struct {
	Hostname string
	Port int
	Suffix string
}

func Use(config *Config, bkd *backend.Backend) *Users {
	conns := newConns(config)
	messages := newMessages(conns)
	conversations := util.NewDummyConversations(messages)
	users := newUsers(conns)

	bkd.Set(conversations, users)

	// TODO: do not return users backend
	return users
}
