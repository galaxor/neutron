package disk

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/emersion/neutron/backend"
)

// Stores private & public keys on disk.
// Public keys are stored in DOMAIN/USERNAME.pub.gpg and private keys are
// in DOMAIN/USERNAME.priv.gpg.
type Keys struct {
	config *Config
	users backend.UsersBackend
}

func (b *Keys) getKeyPath(email string, priv bool) (path string) {
	parts := strings.SplitN(email, "@", 2)

	path = b.config.Directory + "/" + parts[1] + "/" + parts[0]
	if priv {
		path += ".priv"
	} else {
		path += ".pub"
	}
	path += ".gpg"
	return
}

func (b *Keys) getKey(email string, priv bool) (key string, err error) {
	path := b.getKeyPath(email, priv)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	key = string(data)
	return
}

func (b *Keys) GetPublicKey(email string) (string, error) {
	key, err := b.getKey(email, false)
	if err == os.ErrNotExist {
		return "", nil
	}
	return key, err
}

func (b *Keys) GetKeypair(email, password string) (keypair *backend.Keypair, err error) {
	pub, err := b.getKey(email, false)
	if err != nil {
		return
	}

	// TODO: use password to encrypt private key

	priv, err := b.getKey(email, true)
	if err != nil {
		return
	}

	keypair = &backend.Keypair{
		ID: email,
		PublicKey: pub,
		PrivateKey: priv,
	}
	return
}

func (b *Keys) UpdateKeypair(email, password string, keypair *backend.Keypair) error {
	return errors.New("Not yet implemented")
}

func NewKeys(config *Config, users backend.UsersBackend) backend.KeysBackend {
	return &Keys{
		config: config,
		users: users,
	}
}
