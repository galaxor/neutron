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

func Use(config *Config, passwords PasswordsBackend, bkd *backend.Backend) {
	send := New(config, passwords)

	bkd.Set(send)
}
