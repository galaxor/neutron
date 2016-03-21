// Contains a generic interface for backends.
package backend

// A backend takes care of storing all mailbox data.
type Backend interface {
	ContactsBackend
	LabelsBackend
	ConversationsBackend
	SendBackend
	DomainsBackend
	EventsBackend
	SessionsBackend

	// Check if a username is available.
	IsUsernameAvailable(username string) (bool, error)
	// Get a user.
	GetUser(id string) (*User, error)
	// Check if the provided username and password are correct
	Auth(username, password string) (*Session, error)
	// Insert a new user. Returns the newly created user.
	InsertUser(user *User, password string) (*User, error)
	// Update an existing user.
	UpdateUser(update *UserUpdate) error
	// Update a user's password.
	UpdateUserPassword(id, current, new string) error
	// Update a user's private key.
	UpdateKeypair(id, password string, keypair *Keypair) error
	// Delete a user.
	//DeleteUser(id string) error

	// Get a public key for a user.
	GetPublicKey(email string) (string, error)
}
