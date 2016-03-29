package backend

type KeysBackend interface {
	// Get a public key for a user.
	GetPublicKey(email string) (string, error)
	// Get a keypair for a user. Contains public & private key.
	GetKeypair(email, password string) (keypair *Keypair, err error)
	// Update a user's private key.
	UpdateKeypair(email, password string, keypair *Keypair) error
}
