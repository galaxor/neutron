package backend

type KeysBackend interface {
	// Get a public key for a user.
	GetPublicKey(email string) (string, error)
}
