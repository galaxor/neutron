package imap

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type Users struct {
	*conns

	users map[string]*backend.User
}

func (b *Users) GetUser(id string) (user *backend.User, err error) {
	user, ok := b.users[id]
	if !ok {
		err = errors.New("No such user")
	}
	return
}

func (b *Users) Auth(username, password string) (user *backend.User, err error) {
	id := username

	// User already logged in, just checking password
	if client, ok := b.clients[id]; ok {
		if client.password != password {
			err = errors.New("Invalid username or password")
		} else {
			user = b.users[id]
		}
		return
	}

	email, err := b.connect(username, password)
	if err != nil {
		return
	}

	user = &backend.User{
		ID: id,
		Name: username,
		DisplayName: username,
		Addresses: []*backend.Address{
			&backend.Address{
				ID: username,
				Email: email,
				Send: 1,
				Receive: 1,
				Status: 1,
				Type: 1,
			},
		},
	}

	b.users[user.ID] = user

	return
}

func (b *Users) IsUsernameAvailable(username string) (bool, error) {
	return false, errors.New("Cannot check if a username is available with IMAP backend")
}

func (b *Users) InsertUser(u *backend.User, password string) (*backend.User, error) {
	return nil, errors.New("Cannot register new user with IMAP backend")
}

func (b *Users) UpdateUser(update *backend.UserUpdate) error {
	return errors.New("Cannot update user with IMAP backend")
}

func (b *Users) UpdateUserPassword(id, current, new string) error {
	return errors.New("Cannot update user password with IMAP backend")
}

func newUsers(conns *conns) *Users {
	return &Users{
		conns: conns,

		users: map[string]*backend.User{},
	}
}
