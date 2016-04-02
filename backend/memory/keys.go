package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

// A password-protected keypair
type lockedKeypair struct {
	unlocked *backend.Keypair
	password string
}

func (kp *lockedKeypair) GetPublicKey() string {
	return kp.unlocked.PublicKey
}

func (kp *lockedKeypair) Unlock(password string) (*backend.Keypair, error) {
	if kp.password != password {
		return nil, errors.New("Invalid keypair password")
	}
	return kp.unlocked, nil
}

type Keys struct {
	keys map[string]*lockedKeypair
}

func (b *Keys) getLockedKeypair(email string) (*lockedKeypair, error) {
	kp, ok := b.keys[email]
	if !ok {
		return nil, errors.New("No such keypair")
	}
	return kp, nil
}

func (b *Keys) GetPublicKey(email string) (string, error) {
	kp, err := b.getLockedKeypair(email)
	if err != nil {
		return "", nil
	}
	return kp.GetPublicKey(), nil
}

func (b *Keys) GetKeypair(email, password string) (*backend.Keypair, error) {
	kp, err := b.getLockedKeypair(email)
	if err != nil {
		return nil, err
	}

	return kp.Unlock(password)
}

func (b *Keys) UpdateKeypair(email, password string, keypair *backend.Keypair) (updated *backend.Keypair, err error) {
	keypair.ID = email

	b.keys[email] = &lockedKeypair{
		unlocked: keypair,
		password: password,
	}

	updated = keypair
	return
}

func NewKeys() backend.KeysBackend {
	return &Keys{
		keys: map[string]*lockedKeypair{},
	}
}
