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
	item, err := b.getUserData(id)
	if err != nil {
		return
	}

	user = item.user

	if user.EncPrivateKey == "" {
		addr := user.GetMainAddress()
		if addr != nil && len(addr.Keys) > 0 {
			keypair := addr.Keys[0]
			user.PublicKey = keypair.PublicKey
			user.EncPrivateKey = keypair.PrivateKey
		}
	}

	for _, addr := range user.Addresses {
		if addr.DisplayName == "" {
			addr.DisplayName = user.DisplayName
		}
		if len(addr.Keys) > 0 {
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

	// Generate new IDs
	user.ID = generateId()

	for _, addr := range user.Addresses {
		addr.ID = generateId()

		for _, kp := range addr.Keys {
			kp.ID = generateId()
		}
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

	item, err := b.getUserData(updated.ID)
	if err != nil {
		return err
	}

	user := item.user

	if update.DisplayName {
		user.DisplayName = updated.DisplayName
	}
	if update.Signature {
		user.Signature = updated.Signature
	}
	if update.AutoSaveContacts {
		user.AutoSaveContacts = updated.AutoSaveContacts
	}
	if update.ShowImages {
		user.ShowImages = updated.ShowImages
	}

	return nil
}

func checkUserPassword(item *userData, password string) error {
	if item.password != password {
		return errors.New("Invalid password")
	}
	return nil
}

func (b *Backend) UpdateUserPassword(id, current, new string) error {
	item, err := b.getUserData(id)
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

func (b *Backend) UpdateKeypair(id, password string, keypair *backend.Keypair) error {
	item, err := b.getUserData(id)
	if err != nil {
		return err
	}

	err = checkUserPassword(item, password)
	if err != nil {
		return err
	}

	for _, addr := range item.user.Addresses {
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
