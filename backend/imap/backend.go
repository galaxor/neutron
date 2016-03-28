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

func Use(bkd *backend.Backend, config *Config) *conns {
	conns := newConns(config)
	messages := newMessages(conns)
	conversations := util.NewDummyConversations(messages)
	users := newUsers(conns)
	labels := util.NewEventedLabels(newLabels(conns), bkd)

	bkd.Set(messages, conversations, users, labels)

	// TODO: do not return conns backend
	return conns
}
