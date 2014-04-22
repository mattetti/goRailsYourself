package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func (crypt *MessageEncryptor) aesCbcEncrypt(value interface{}) (string, error) {
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

	// Set a default serializer if not already set
	if crypt.Serializer == nil {
		crypt.Serializer = JsonMsgSerializer{}
	}
	splaintext, err := crypt.Serializer.Serialize(value)
	if err != nil {
		return "", err
	}
	plaintext := []byte(splaintext)

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. See
	// http://tools.ietf.org/html/rfc5652#section-6.3
	plaintext = PKCS7Pad(plaintext)

	// The IV needs to be unique, but not secure, it is included in the
	// cypher text.
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// generate the cipher text
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// base64 the cipher text + the iv and join by "--"
	output := base64.StdEncoding.EncodeToString(ciphertext) + "--" + base64.StdEncoding.EncodeToString(iv)
	return output, nil
}

func (crypt *MessageEncryptor) aesCbcDecrypt(encryptedMsg string, target interface{}) error {
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

	// split the msg and decode each part
	splitMsg := strings.Split(encryptedMsg, "--")
	if len(splitMsg) != 2 {
		return errors.New("bad data (--)")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(splitMsg[0])
	if err != nil {
		return err
	}
	iv, err := base64.StdEncoding.DecodeString(splitMsg[1])
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return errors.New("bad data, ciphertext too short")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return errors.New("bad data, ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	unPaddedCiphertext := PKCS7Unpad(ciphertext)

	// In some cases, Rails sends us messages padded with 0x10 (while this package only pads with 0x01-0x0f).
	// For now, we handle this case here when the Serializer is JSON (so we know that 0x10 is actually a padding
	// and not valid data - because this is an invalid json character).
	if _, ok := crypt.Serializer.(JsonMsgSerializer); ok {
		unPaddedCiphertext = bytes.TrimRight(unPaddedCiphertext, "\x10")
	}

	return crypt.Serializer.Unserialize(string(unPaddedCiphertext), target)
}
