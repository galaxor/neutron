package smtp

import (
	"net/smtp"
	"strconv"

	"github.com/emersion/neutron/backend"
	"github.com/emersion/neutron/backend/util/textproto"
)

type PasswordsBackend interface {
	GetPassword(user string) (string, error)
}

type Config struct {
	Hostname string
	Port int
	Suffix string
}

type SendBackend struct {
	PasswordsBackend

	config *Config
}

func (b *SendBackend) SendMessagePackage(user string, msg *backend.OutgoingMessage) error {
	password, err := b.GetPassword(user)
	if err != nil {
		return err
	}

	host := b.config.Hostname
	if b.config.Port > 0 {
		host += ":" + strconv.Itoa(b.config.Port)
	}

	auth := smtp.PlainAuth("", user + b.config.Suffix, password, b.config.Hostname)
	recipients := []string{msg.MessagePackage.Address}
	mail := textproto.FormatOutgoingMessage(msg)

	err = smtp.SendMail(host, auth, msg.Sender.Address, recipients, []byte(mail))
	if err != nil {
		return err
	}

	return nil
}

func New(config *Config, passwords PasswordsBackend) backend.SendBackend {
	return &SendBackend{
		PasswordsBackend: passwords,

		config: config,
	}
}
