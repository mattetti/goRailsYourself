package crypto

import (
	"crypto/rand"
	"io"
)

type MsgSerializer interface {
	Serialize(v interface{}) (string, error)
	Unserialize(data string, v interface{}) error
}

// Generates a random key of the passed length.
// As a reminder, for AES keys of length 16, 24, or 32 bytes are expected for AES-128, AES-192, or AES-256.
func GenerateRandomKey(strength int) []byte {
	k := make([]byte, strength)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
