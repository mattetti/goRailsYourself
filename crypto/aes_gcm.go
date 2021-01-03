package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func (crypt *MessageEncryptor) aesGCMEncrypt(value interface{}) (string, error) {
	// TODO: check the crypt is properly initiated
	k := crypt.Key
	// The longest accepted key is 32 byte long,
	// instead of rejecting a long key, we truncate it.
	// This is how openssl in Ruby works.
	if len(k) > 32 {
		k = crypt.Key[:32]
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Set a default serializer if not already set
	if crypt.Serializer == nil {
		crypt.Serializer = JsonMsgSerializer{}
	}
	splaintext, err := crypt.Serializer.Serialize(value)
	if err != nil {
		return "", err
	}
	plaintext := []byte(splaintext)

	iv := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)

	// Rails stores the GCM auth tag separately from the encrypted data,
	// unlike the cipher package, so a little munging is required.
	// Luckily aesgcm.Overhead() is the tag size (which is 16).
	tagStart := len(ciphertext) - aesgcm.Overhead()
	tag := ciphertext[tagStart:]
	enc := ciphertext[:tagStart]

	vectors := [][]byte{enc, iv, tag}
	for i, vec := range vectors {
		dst := make([]byte, base64.StdEncoding.EncodedLen(len(vec)))
		base64.StdEncoding.Encode(dst, vec)
		vectors[i] = dst
	}

	output := string(bytes.Join(vectors, []byte("--")))
	return output, nil
}

func (crypt *MessageEncryptor) aesGCMDecrypt(encryptedMsg string, target interface{}) error {
	k := crypt.Key
	// The longest accepted key is 32 byte long,
	// instead of rejecting a long key, we truncate it.
	// This is how openssl in Ruby works.
	if len(k) > 32 {
		k = crypt.Key[:32]
	}

	block, err := aes.NewCipher(k)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	vectors := bytes.SplitN([]byte(encryptedMsg), []byte("--"), 3)
	if len(vectors) != 3 {
		return fmt.Errorf("missing vectors, want 3, got %d", len(vectors))
	}
	for i, vec := range vectors {
		dst := make([]byte, base64.StdEncoding.DecodedLen(len(vec)))
		n, err := base64.StdEncoding.Decode(dst, vec)
		if err != nil {
			return fmt.Errorf("bad base64 encoding")
		}
		vectors[i] = dst[:n]
	}

	enc := vectors[0]
	// Rails splits the auth tag into a separate vector, which is unnecessary really, but fine.
	enc = append(enc, vectors[2]...)
	nonce := vectors[1]

	plain, err := aesgcm.Open(nil, nonce, enc, nil)
	if err != nil {
		return err
	}

	return crypt.Serializer.Unserialize(string(plain), target)
}
