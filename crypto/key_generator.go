package crypto

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// KeyGenerator is a simple wrapper around a PBKDF2 implementation.
// It can be used to derive a number of keys for various purposes from a given secret.
// This lets applications have a single secure secret, but avoid reusing that
// key in multiple incompatible contexts.
type KeyGenerator struct {
	Secret     string
	Iterations int
	cache      map[string][]byte
}

// CacheGenerate() write through cache used to save generated keys.
func (g *KeyGenerator) CacheGenerate(salt []byte, keySize int) []byte {
	key := fmt.Sprintf("%s%d", salt, keySize)
	if g.cache == nil {
		g.cache = map[string][]byte{}
	}
	if g.cache[key] == nil {
		g.cache[key] = g.Generate(salt, keySize)
	}
	return g.cache[key]
}

// Generates a derived key based on a salt. rails default key size is 64.
func (g *KeyGenerator) Generate(salt []byte, keySize int) []byte {
	// set a default
	if g.Iterations == 0 {
		g.Iterations = 1000 // rails 4 default when setting the session.
	}
	return pbkdf2.Key([]byte(g.Secret), salt, g.Iterations, keySize, sha1.New)
}
