package smtp

import (
	"github.com/emersion/neutron/backend"
)

type PasswordsBackend interface {
	GetPassword(user string) (string, error)
}

type Config struct {
	Hostname string
	Port int
	Suffix string
}

func Use(bkd *backend.Backend, config *Config, passwords PasswordsBackend) {
	send := New(config, passwords)

	bkd.Set(send)
}
