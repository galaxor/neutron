package memory

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

const idLength = 64

func generateId() string {
	b := make([]byte, idLength)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	id := base64.StdEncoding.EncodeToString(b)
	return strings.Replace(id, "/", "_", -1)
}
