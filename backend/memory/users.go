package memory

import (
	"errors"
	"io/ioutil"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) getKeypair(id string) (keypair *backend.Keypair, err error) {
	var pub []byte
	pub, err = ioutil.ReadFile("data/public.key")
	if err != nil {
		return
	}

	var priv []byte
	priv, err = ioutil.ReadFile("data/private.key")
	if err != nil {
		return
	}

	keypair = &backend.Keypair{
		ID: "keypair_id",
		PublicKey: string(pub),
		PrivateKey: string(priv),
	}
	return
}

func (b *Backend) IsUsernameAvailable(username string) (bool, error) {
	for _, d := range b.data {
		if d.user.Name == username {
			return false, nil
		}
	}

	return true, nil
}

func (b *Backend) GetUser(id string) (user *backend.User, err error) {
	item, ok := b.data[id]
	if !ok {
		err = errors.New("No such user")
		return
	}

	user = item.user

	keypair, err := b.getKeypair(id)
	if err != nil {
		return
	}

	user.PublicKey = keypair.PublicKey
	user.EncPrivateKey = keypair.PrivateKey

	user.Addresses = []*backend.Address{
		&backend.Address{
			ID: "address_id",
			DomainID: "domain_id",
			Email: "neutron@example.org",
			Send: 1,
			Receive: 1,
			DisplayName: "Neutron",
			Keys: []*backend.Keypair{keypair},
		},
	}

	return
}

func (b *Backend) Auth(username, password string) (user *backend.User, err error) {
	for id, item := range b.data {
		if item.user.Name == username && item.password == password {
			user, err = b.GetUser(id)
			return
		}
	}

	err = errors.New("Invalid username and password combination")
	return
}

func (b *Backend) InsertUser(user *backend.User, password string) (*backend.User, error) {
	user.ID = "user_id" // TODO

	b.data[user.ID] = &userData{
		user: user,
		password: password,
	}

	return b.GetUser(user.ID)
}
