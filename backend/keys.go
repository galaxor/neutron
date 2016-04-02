package backend

type KeysBackend interface {
	// Get a public key for a user.
	// If no key is available, an empty string and no error must be returned.
	GetPublicKey(email string) (string, error)
	// Get a keypair for a user. Contains public & private key.
	GetKeypair(email string) (*Keypair, error)
	// Create a new keypair.
	InsertKeypair(email string, keypair *Keypair) (*Keypair, error)
	// Update a user's private key.
	// PublicKey must be updated only if it isn't empty.
	UpdateKeypair(email string, keypair *Keypair) (*Keypair, error)
}
