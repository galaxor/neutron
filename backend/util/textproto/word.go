package textproto

import (
	"mime"
)

// Decode a RFC2047 word
func DecodeWord(word string) string {
	dec := new(mime.WordDecoder) // TODO: do not create one decoder per word
	decoded, err := dec.DecodeHeader(word)
	if err == nil {
		return decoded
	}
	return word
}
