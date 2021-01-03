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
// Different kind of ciphers are supported:
//  - aes-cbc - Rails' default until 5.2, requires a verifier
//  - aes-256-gcm - Rails 5.2+ default, ignores verifier.
//
// Note: The old Rails default serializer, Marshal is neither safe or
// portable across langauges, use the JSON serializer.
type MessageEncryptor struct {
	Key []byte
	// optional property used to automatically set the
	// verifier if not already set.
	SignKey    []byte
	Cipher     string
	Verifier   *MessageVerifier
	Serializer MsgSerializer
}

func (crypt *MessageEncryptor) withVerifier() bool {
	switch crypt.Cipher {
	case "aes-256-gcm":
		return false
	}
	return true
}

// EncryptAndSign performs encryption with authentication, or encryption
// followed by signing, depending on the selected cipher mode. message can be
// any serializable type (string, struct, map, etc).
// Note that even if you can just Encrypt() in most cases you shouldn't use it
// directly and instead use this method.
// For aes-cbc mode, encryption alone is neither signed or authenticated, and is
// subject to padding oracle attacks.
// Reference: http://www.limited-entropy.com/padding-oracle-attacks.
// The output string can be converted back using DecryptAndVerify() and is
// encoded using base64.
func (crypt *MessageEncryptor) EncryptAndSign(value interface{}) (string, error) {
	if crypt == nil {
		return "", errors.New("can't call EncryptAndSign on a nil *MessageEncryptor")
	}

	if !crypt.withVerifier() {
		return crypt.Encrypt(value)
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

// DecryptAndVerify decrypts and either authenticates or verifies the signature
// of a message, depending on the selected cipher mode. Messages need to be
// either signed or authenticated (GCM) on top of being encrypted in order to
// avoid padding attacks. Reference: http://www.limited-entropy.com/padding-oracle-attacks.
// The serializer will populate the pointer you are passing as second argument.
func (crypt *MessageEncryptor) DecryptAndVerify(msg string, target interface{}) error {

	if !crypt.withVerifier() {
		return crypt.Decrypt(msg, target)
	}

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

// Encrypt encrypts a message using the set cipher and the secret.
// The returned value is a base 64 encoded string of the encrypted data + IV joined by "--".
// An encrypted message isn't safe unless it's signed!
func (crypt *MessageEncryptor) Encrypt(value interface{}) (string, error) {
	switch crypt.Cipher {
	case "aes-cbc":
		return crypt.aesCbcEncrypt(value)
	case "aes-256-gcm":
		return crypt.aesGCMEncrypt(value)
	case "":
		// using a default if not set
		return crypt.aesCbcEncrypt(value)
	}
	return "", errors.New("cipher not set or not supported")
}

// Decrypt decrypts a message using the set cipher and the secret.
// The passed value is expected to be a base 64 encoded string of the encrypted data + IV joined by "--"
func (crypt *MessageEncryptor) Decrypt(value string, target interface{}) error {
	if crypt.Serializer == nil {
		crypt.Serializer = JsonMsgSerializer{}
	}
	switch crypt.Cipher {
	case "aes-cbc":
		return crypt.aesCbcDecrypt(value, target)
	case "aes-256-gcm":
		return crypt.aesGCMDecrypt(value, target)
	case "":
		// using a default if not set
		return crypt.aesCbcDecrypt(value, target)
	}
	return errors.New("cipher not set or not supported")
}
