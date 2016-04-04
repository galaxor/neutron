// Backend utilities.
package util

import (
	"crypto/rand"
	"encoding/base64"
)

const idLength = 64

// Generate a random identifier.
func GenerateId() string {
	b := make([]byte, idLength)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
