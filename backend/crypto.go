package backend

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

type Keypair struct {
	ID string
	PublicKey string
	PrivateKey string
}

func (kp *Keypair) EncryptToSelf(data string) (encrypted string, err error) {
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
	armorWriter, err := armor.Encode(&tokenBuffer, "PGP MESSAGE", map[string]string{})
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
