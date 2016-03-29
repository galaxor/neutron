package disk

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type Config struct {
	Directory string
}

// Stores private & public keys on disk.
// Public keys are stored in DOMAIN/USERNAME.pub.gpg and private keys are
// in DOMAIN/USERNAME.priv.gpg.
type Keys struct {
	config *Config
}

func (b *Keys) GetPublicKey(email string) (string, error) {
	return "", errors.New("Not yet implemented")
}

func (b *Keys) UpdateKeypair(id, password string, keypair *backend.Keypair) error {
	return errors.New("Not yet implemented")
}

func NewKeys(config *Config, users backend.UsersBackend) backend.KeysBackend {
	return &Keys{
		config: config,
	}
}
