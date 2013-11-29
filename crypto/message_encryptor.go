package crypto

import (
	"crypto/sha1"
	"errors"
)

//
// MessageEncryptor is a simple way to encrypt values which get stored
// somewhere you don't trust.
//
// The cipher text and initialization vector are base64 encoded and returned
// to you.
//
// This can be used in situations similar to the MessageVerifier, but
// where you don't want users to be able to determine the value of the payload.
//
// Different kind of ciphers will be supported, currently only Rails' default aes-cbc
// is supported.
// Note that as I'm writing this library, Rails default serializer is Ruby's Marshal
// which is not safe and cross language. You need to switch the serializer to JSON or another
// safer/cross language format to share encrypted messages between Ruby and Go.

// Ruby/Rails 4 usage:
// require 'active_support'
// require 'json'
// salt  = SecureRandom.random_bytes(64)
// key   = ActiveSupport::KeyGenerator.new('password').generate_key(salt)
// crypt = ActiveSupport::MessageEncryptor.new(key, "this is the sign secret", serializer: JSON)
// encrypted_data = crypt.encrypt_and_sign('my secret data')              # => "emsxbm5HcVJWRmhZTzNPTEFjTERHUjJjbmpIWXF5UzNITWhMem5sUnNZRT0tLVVCak1GeDFrSHVxaGFyeVpqRlVLNHc9PQ==--789d60509d8b441a24600bbf48af47d3eff386b5"
// crypt.decrypt_and_verify(encrypted_data)                               # => "my secret data"
type MessageEncryptor struct {
	Key []byte
	// optionial property used to automatically set the
	// verifier if not already set.
	SignKey    []byte
	cipher     string
	verifier   *MessageVerifier
	serializer MsgSerializer
}

// Encrypt and sign a message (string, struct, anything that can be safely serialized/serialized).
// Note that even if you can just Encrypt()
// in most cases you shouldn't use it directly and instead use this method.
// The reason is that we need to sign the message in order to avoid
// padding attacks.
// Reference: http://www.limited-entropy.com/padding-oracle-attacks.
// The output string can be converted back using DecryptAndVerify()
// and is encoded using base64.
func (crypt *MessageEncryptor) EncryptAndSign(value interface{}) (string, error) {
	// Set a default verifier if a signature key was given instead of setting the verifier directly.
	if crypt.verifier == nil && crypt.SignKey != nil {
		crypt.verifier = &MessageVerifier{
			secret:     crypt.SignKey,
			hasher:     sha1.New,
			serializer: NullMsgSerializer{},
		}
	}
	if crypt.verifier == nil {
		return "", errors.New("Verifier and/or signature key not set: ")
	}
	vvalid, err := crypt.verifier.IsValid()
	if !vvalid {
		return "", errors.New("Verifier not properly set: " + err.Error())
	}
	encryptedMsg, err := crypt.Encrypt(value)
	if err != nil {
		return "", err
	}
	return crypt.verifier.Generate(encryptedMsg)
}

// Decrypt and verify a message. Messages need to be signed on top of being encrypted in order to
// avoid padding attacks. Reference: http://www.limited-entropy.com/padding-oracle-attacks.
// The serializer will populate the pointer you are passing as second argument.
func (crypt *MessageEncryptor) DecryptAndVerify(msg string, target interface{}) error {
	// Set a default verifier if a signature key was given instead of setting the verifier directly.
	if crypt.verifier == nil && crypt.SignKey != nil {
		crypt.verifier = &MessageVerifier{
			secret:     crypt.SignKey,
			hasher:     sha1.New,
			serializer: NullMsgSerializer{},
		}
	}
	var base64Msg string
	// verify the data and get the encoded data out.
	err := crypt.verifier.Verify(msg, &base64Msg)
	if err != nil {
		return errors.New("Verification failed: " + err.Error())
	}
	return crypt.Decrypt(base64Msg, target)
}

// Encrypt() encrypts a message using the set cipher and the secret.
// The returned value is a base 64 encoded string of the encrypted data + IV joined by "--".
// An encrypted message isn't safe unless it's signed!
func (crypt *MessageEncryptor) Encrypt(value interface{}) (string, error) {
	switch crypt.cipher {
	case "aes-cbc":
		return crypt.aesCbcEncrypt(value)
	case "":
		// using a default if not set
		return crypt.aesCbcEncrypt(value)
	}
	return "", errors.New("cipher not set or not supported")
}

// decrypt() decrypts a message using the set cipher and the secret.
// The passed value is expected to be a base 64 encoded string of the encrypted data + IV joined by "--"
func (crypt *MessageEncryptor) Decrypt(value string, target interface{}) error {
	if crypt.serializer == nil {
		crypt.serializer = JsonMsgSerializer{}
	}
	switch crypt.cipher {
	case "aes-cbc":
		return crypt.aesCbcDecrypt(value, target)
	case "":
		// using a default if not set
		return crypt.aesCbcDecrypt(value, target)
	}
	return errors.New("cipher not set or not supported")
}
