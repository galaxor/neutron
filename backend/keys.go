package backend

type KeysBackend interface {
	// Get a public key for a user.
	// If no key is available, an empty string and no error must be returned.
	GetPublicKey(email string) (string, error)
	// Get a keypair for a user. Contains public & private key.
	GetKeypair(email, password string) (*Keypair, error)
	// Update a user's private key.
	UpdateKeypair(email, password string, keypair *Keypair) (*Keypair, error)
}
