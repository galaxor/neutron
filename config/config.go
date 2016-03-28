package config

import (
	"github.com/emersion/neutron/backend/imap"
	"github.com/emersion/neutron/backend/smtp"
)

// Configuration for all backends.
// Backends omitted or set to null won't be activated.
type Config struct {
	// Activate the memory backend.
	Memory *MemoryConfig
	// IMAP config.
	Imap *ImapConfig
	// SMTP config.
	Smtp *SmtpConfig
}

type BackendConfig struct {
	Enabled bool
}

type MemoryConfig struct {
	BackendConfig
	Populate bool
}

type ImapConfig struct {
	BackendConfig
	imap.Config
}

type SmtpConfig struct {
	BackendConfig
	smtp.Config
}
