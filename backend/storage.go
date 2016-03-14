package backend

import (
	"errors"
	"io/ioutil"
)

type User struct {
	Uid string
	Username string
	PrivateKey string
}

func Login(username, password string) (user *User, err error) {
	if username != "neutron" || password != "neutron" {
		err = errors.New("Invalid username and password combination")
		return
	}

	user = &User{
		Uid: "0000000000000000000000000000000000000000",
		Username: username,
	}

	var priv []byte
	priv, err = ioutil.ReadFile("data/private.key")
	if err != nil {
		return
	}

	user.PrivateKey = string(priv)

	return
}
