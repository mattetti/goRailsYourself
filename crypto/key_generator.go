package crypto

import (
  "crypto/sha1"
  "code.google.com/p/go.crypto/pbkdf2"
)

// KeyGenerator is a simple wrapper around a PBKDF2 implementation.
// It can be used to derive a number of keys for various purposes from a given secret.
// This lets applications have a single secure secret, but avoid reusing that
// key in multiple incompatible contexts.
type KeyGenerator struct {
  Secret string
  Iterations int
}

func (g *KeyGenerator) Generate(salt []byte, keySize int) []byte {
  // set a default
  if g.Iterations == 0 {
    g.Iterations = 64
  }
  return pbkdf2.Key([]byte(g.Secret), salt, 4096, keySize, sha1.New)
}
