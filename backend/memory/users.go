package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type Users struct {
	users map[string]*user
}

type user struct {
	*backend.User
	password string
}

func (b *Users) IsUsernameAvailable(username string) (bool, error) {
	for _, d := range b.users {
		if d.Name == username {
			return false, nil
		}
	}

	return true, nil
}

func (b *Users) getUser(id string) (*user, error) {
	user, ok := b.users[id]
	if !ok {
		return nil, errors.New("No such user")
	}
	return user, nil
}

func (b *Users) GetUser(id string) (user *backend.User, err error) {
	item, err := b.getUser(id)
	if err != nil {
		return
	}

	user = item.User

	if user.EncPrivateKey == "" {
		addr := user.GetMainAddress()
		if addr != nil && len(addr.Keys) > 0 {
			keypair := addr.Keys[0]
			user.PublicKey = keypair.PublicKey
			user.EncPrivateKey = keypair.PrivateKey
		}
	}

	return
}

func (b *Users) Auth(username, password string) (user *backend.User, err error) {
	for id, item := range b.users {
		if item.Name == username && item.password == password {
			return b.GetUser(id)
		}
	}

	err = errors.New("Invalid username and password combination")
	return
}

func (b *Users) InsertUser(u *backend.User, password string) (*backend.User, error) {
	available, err := b.IsUsernameAvailable(u.Name)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("Username already taken")
	}

	// Generate new IDs
	u.ID = generateId()

	for _, addr := range u.Addresses {
		addr.ID = generateId()

		for _, kp := range addr.Keys {
			kp.ID = generateId()
		}
	}

	// Insert new user
	b.users[u.ID] = &user{
		User: u,
		password: password,
	}

	return b.GetUser(u.ID)
}

func (b *Users) UpdateUser(update *backend.UserUpdate) error {
	item, err := b.getUser(update.User.ID)
	if err != nil {
		return err
	}

	user := item.User
	update.Apply(user)
	return nil
}

func checkUserPassword(item *user, password string) error {
	if item.password != password {
		return errors.New("Invalid password")
	}
	return nil
}

func (b *Users) UpdateUserPassword(id, current, new string) error {
	item, err := b.getUser(id)
	if err != nil {
		return err
	}

	err = checkUserPassword(item, current)
	if err != nil {
		return err
	}

	item.password = new
	return nil
}

func (b *Users) UpdateKeypair(id, password string, keypair *backend.Keypair) error {
	item, err := b.getUser(id)
	if err != nil {
		return err
	}

	err = checkUserPassword(item, password)
	if err != nil {
		return err
	}

	for _, addr := range item.User.Addresses {
		for _, kp := range addr.Keys {
			if kp.ID == keypair.ID {
				kp.PrivateKey = keypair.PrivateKey
				kp.PublicKey = "" // Public key is now outdated
				return nil
			}
		}
	}

	return errors.New("Key not found")
}

func (b *Users) GetPublicKey(email string) (string, error) {
	for _, item := range b.users {
		for _, address := range item.User.Addresses {
			if address.Email == email && len(address.Keys) > 0 {
				return address.Keys[0].PublicKey, nil
			}
		}
	}
	return "", nil
}

func NewUsers() *Users {
	return &Users{
		users: map[string]*user{},
	}
}
