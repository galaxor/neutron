package backend

import (
	"bytes"
	"errors"
	"strings"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

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
}

// Encrypt a message to the keypair's owner.
func (kp *Keypair) Encrypt(data string) (encrypted string, err error) {
	entitiesList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(kp.PrivateKey))
	if err != nil {
		return
	}
	if len(entitiesList) == 0 {
		err = errors.New("Key ring does not contain any key")
		return
	}

	entity := entitiesList[0]

	var tokenBuffer bytes.Buffer
	armorWriter, err := ArmorMessage(&tokenBuffer)
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

	encrypted = tokenBuffer.String()
	return
}

// Check if a string contains an encrypted message.
func IsEncrypted(data string) bool {
	return strings.Contains(data, "-----BEGIN " + PgpMessageType + "-----")
}
