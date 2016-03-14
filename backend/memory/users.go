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

func (b *Backend) GetUser(id string) (user *backend.User, err error) {
	user = &backend.User{
		ID: id,
		Name: "neutron",
		DisplayName: "Neutron",
		NotificationEmail: "neutron@example.org",
	}

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
	if username != "neutron" || password != "neutron" {
		err = errors.New("Invalid username and password combination")
		return
	}

	user, err = b.GetUser("user_id")
	return
}
