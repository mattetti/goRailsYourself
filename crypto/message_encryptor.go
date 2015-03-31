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
type MessageEncryptor struct {
	Key []byte
	// optional property used to automatically set the
	// verifier if not already set.
	SignKey    []byte
	Cipher     string
	Verifier   *MessageVerifier
	Serializer MsgSerializer
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
	if crypt == nil {
		return "", errors.New("can't call EncryptAndSign on a nil *MessageEncryptor")
	}
	// Set a default verifier if a signature key was given instead of setting the verifier directly.
	if crypt.Verifier == nil && crypt.SignKey != nil {
		crypt.Verifier = &MessageVerifier{
			Secret:     crypt.SignKey,
			Hasher:     sha1.New,
			Serializer: NullMsgSerializer{},
		}
	}
	if crypt.Verifier == nil {
		return "", errors.New("Verifier and/or signature key not set: ")
	}
	vvalid, err := crypt.Verifier.IsValid()
	if !vvalid {
		return "", errors.New("Verifier not properly set: " + err.Error())
	}
	encryptedMsg, err := crypt.Encrypt(value)
	if err != nil {
		return "", err
	}
	return crypt.Verifier.Generate(encryptedMsg)
}

// Decrypt and verify a message. Messages need to be signed on top of being encrypted in order to
// avoid padding attacks. Reference: http://www.limited-entropy.com/padding-oracle-attacks.
// The serializer will populate the pointer you are passing as second argument.
func (crypt *MessageEncryptor) DecryptAndVerify(msg string, target interface{}) error {
	// Set a default verifier if a signature key was given instead of setting the verifier directly.
	if crypt.Verifier == nil && crypt.SignKey != nil {
		crypt.Verifier = &MessageVerifier{
			Secret:     crypt.SignKey,
			Hasher:     sha1.New,
			Serializer: NullMsgSerializer{},
		}
	}
	var base64Msg string
	// verify the data and get the encoded data out.
	err := crypt.Verifier.Verify(msg, &base64Msg)
	if err != nil {
		return errors.New("Verification failed: " + err.Error())
	}
	return crypt.Decrypt(base64Msg, target)
}

// Encrypt() encrypts a message using the set cipher and the secret.
// The returned value is a base 64 encoded string of the encrypted data + IV joined by "--".
// An encrypted message isn't safe unless it's signed!
func (crypt *MessageEncryptor) Encrypt(value interface{}) (string, error) {
	switch crypt.Cipher {
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
	if crypt.Serializer == nil {
		crypt.Serializer = JsonMsgSerializer{}
	}
	switch crypt.Cipher {
	case "aes-cbc":
		return crypt.aesCbcDecrypt(value, target)
	case "":
		// using a default if not set
		return crypt.aesCbcDecrypt(value, target)
	}
	return errors.New("cipher not set or not supported")
}
