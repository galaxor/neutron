package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

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

	keypair := &backend.Keypair{
		PublicKey: user.PublicKey,
		PrivateKey: user.EncPrivateKey,
	}

	for _, addr := range user.Addresses {
		if addr.DisplayName == "" {
			addr.DisplayName = user.DisplayName
		}
		if len(addr.Keys) == 0 {
			addr.Keys = []*backend.Keypair{keypair}
			addr.HasKeys = 1
		}
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
	available, err := b.IsUsernameAvailable(user.Name)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("Username already taken")
	}

	user.ID = generateId()

	user.Addresses = []*backend.Address{
		&backend.Address{
			ID: generateId(),
			DomainID: "domain_id", // TODO
			Email: user.Name + "@example.org", // TODO
			Send: 1,
			Receive: 1,
			Status: 1,
			Type: 1,
		},
	}

	// Insert new user
	b.data[user.ID] = &userData{
		user: user,
		password: password,
	}

	return b.GetUser(user.ID)
}

func (b *Backend) UpdateUser(update *backend.UserUpdate) error {
	updated := update.User

	item, ok := b.data[updated.ID]
	if !ok {
		return errors.New("No such user")
	}

	user := item.user

	if update.DisplayName {
		user.DisplayName = updated.DisplayName
	}
	if update.Signature {
		user.Signature = updated.Signature
	}

	return nil
}

func (b *Backend) GetPublicKey(email string) (string, error) {
	for _, data := range b.data {
		for _, address := range data.user.Addresses {
			if address.Email == email && len(address.Keys) > 0 {
				return address.Keys[0].PublicKey, nil
			}
		}
	}
	return "", nil
}
