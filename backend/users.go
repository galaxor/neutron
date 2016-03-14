package backend

import (
	"errors"
	"io/ioutil"
)

type User struct {
	ID string
	Name string
	NotificationEmail string
	Signature string
	NumMessagePerPage int
	UsedSpace int
	Notify int
	AutoSaveContacts int
	Language string
	LogAuth int
	ComposerMode int
	MessageButtons int
	ShowImages int
	ViewMode int
	ViewLayout int
	SwipeLeft int
	SwipeRight int
	Theme string
	Currency string
	Credit int
	DisplayName string
	MaxSpace int
	MaxUpload int
	Role int
	Private int
	Subscribed int
	Deliquent int
	Addresses []*Address
	PublicKey string
	EncPrivateKey string
}

type Address struct {
	ID string
	DomainID string
	Email string
	Send int
	Receive int
	Status int
	Type int
	DisplayName string
	Signature string
	HashKeys int
	Keys []*Keypair
}

type Keypair struct {
	ID string
	PublicKey string
	PrivateKey string
}

func GetKeypair(id string) (keypair *Keypair, err error) {
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

	keypair = &Keypair{
		ID: "keypair_id",
		PublicKey: string(pub),
		PrivateKey: string(priv),
	}
	return
}

func Get(id string) (user *User, err error) {
	user = &User{
		ID: id,
		Name: "neutron",
		DisplayName: "Neutron",
		NotificationEmail: "neutron@example.org",
	}

	keypair, err := GetKeypair(id)
	if err != nil {
		return
	}

	user.PublicKey = keypair.PublicKey
	user.EncPrivateKey = keypair.PrivateKey

	user.Addresses = []*Address{
		&Address{
			ID: "address_id",
			DomainID: "domain_id",
			Email: "neutron@example.org",
			Send: 1,
			Receive: 1,
			DisplayName: "Neutron",
			Keys: []*Keypair{keypair},
		},
	}

	return
}

func Auth(username, password string) (user *User, err error) {
	if username != "neutron" || password != "neutron" {
		err = errors.New("Invalid username and password combination")
		return
	}

	user, err = Get("user_id")
	return
}
