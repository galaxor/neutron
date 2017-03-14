package backend

import (
	"bytes"
	"errors"
	"strings"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

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


const PgpMessageType = "PGP MESSAGE"

// Encode a PGP message armor.
func ArmorMessage(w io.Writer) (io.WriteCloser, error) {
	return armor.Encode(w, PgpMessageType, map[string]string{})
}

// A keypair contains a private and a public key.
type Keypair struct {
	ID string
	PublicKey string
	PrivateKey string
	Fingerprint string // TODO: populate this field
}

func (kp *Keypair) getPrivateKey() (entity *openpgp.Entity, err error) {
	entitiesList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(kp.PrivateKey))
	if err != nil {
		return
	}

	if len(entitiesList) == 0 {
		err = errors.New("Key ring does not contain any key")
		return
	}

	entity = entitiesList[0]
	return
}

// Encrypt a message to the keypair's owner.
func (kp *Keypair) Encrypt(data string) (encrypted string, err error) {
	entity, err := kp.getPrivateKey()
	if err != nil {
		return
	}

	var b bytes.Buffer
	armorWriter, err := ArmorMessage(&b)
	if err != nil {
		return
	}

	w, err := openpgp.Encrypt(armorWriter, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return
	}

	w.Write([]byte(data))
	w.Close()

	armorWriter.Close()

	encrypted = b.String()
	return
}

// Read public key from private key
func (kp *Keypair) readPublicKey() (err error) {
	entity, err := kp.getPrivateKey()
	if err != nil {
		return
	}

	var b bytes.Buffer
	w, err := armor.Encode(&b, openpgp.PublicKeyType, nil)
	if err != nil {
		return
	}

	err = entity.Serialize(w)
	if err != nil {
		return
	}
	w.Close()

	kp.PublicKey = b.String()
	return
}

// Create a new keypair.
func NewKeypair(pub, priv string) *Keypair {
	kp := &Keypair{
		PublicKey: pub,
		PrivateKey: priv,
	}

	if kp.PublicKey == "" {
		err := kp.readPublicKey()
		if err != nil {
			panic(err)
		}
	}

	return kp
}

// Check if a string contains an encrypted message.
func IsEncrypted(data string) bool {
	return strings.Contains(data, "-----BEGIN " + PgpMessageType + "-----")
}
