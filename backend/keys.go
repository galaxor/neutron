package backend

type KeysBackend interface {
	// Get a public key for a user.
	GetPublicKey(email string) (string, error)
	// Update a user's private key.
	UpdateKeypair(id, password string, keypair *Keypair) error
}
