package crypto

import (
	"crypto/sha1"
	"fmt"
	. "github.com/franela/goblin"
	"strings"
	"testing"
)

func TestMessageEncryptorDefaultSettings(t *testing.T) {
	g := Goblin(t)

	g.Describe("MessageEncryptor with default settings", func() {
		k := GenerateRandomKey(32)
		signKey := []byte("this is a secret!")
		e := MessageEncryptor{Key: k, SignKey: signKey}
		g.It("can round trip an encoded/unsigned string", func() {
			msg, err := e.Encrypt("my secret data")
			g.Assert(err).Eql(nil)
			var newMsg string
			err = e.Decrypt(msg, &newMsg)
			g.Assert(err).Eql(nil)
			g.Assert(newMsg).Eql("my secret data")
		})
		g.It("can round trip an encoded/signed string", func() {
			msg, err := e.EncryptAndSign("my secret data")
			g.Assert(err).Eql(nil)
			var newMsg string
			err = e.DecryptAndVerify(msg, &newMsg)
			g.Assert(err).Eql(nil)
			g.Assert(newMsg).Eql("my secret data")
		})

	})

}

func TestMessageEncryptor(t *testing.T) {
	g := Goblin(t)

	g.Describe("MessageEncryptor properly setup using aes cbc", func() {
		newCrypt := func() MessageEncryptor {
			return MessageEncryptor{Key: GenerateRandomKey(32),
				Cipher: "aes-cbc",
				Verifier: &MessageVerifier{
					Secret:     []byte("signature secret!"),
					Hasher:     sha1.New,
					Serializer: NullMsgSerializer{},
				},
				Serializer: JsonMsgSerializer{},
			}
		}

		g.It("can encrypt/decrypt an unsigned string", func() {
			e := newCrypt()
			msg, err := e.Encrypt("my secret data")
			g.Assert(err).Eql(nil)
			splitMsg := strings.Split(msg, "--")
			g.Assert(len(splitMsg)).Eql(2)
			//encryptedMsg, iv := splitMsg[0], splitMsg[1]
			var newMsg string
			err = e.Decrypt(msg, &newMsg)
			g.Assert(err).Eql(nil)
			g.Assert(newMsg).Eql("my secret data")
		})

		g.It("can encrypt/decrypt an unsigned struct", func() {
			type Person struct {
				Id        int    `json:"id"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Age       int    `json:"age"`
			}
			data := Person{Id: 12, FirstName: "John", LastName: "Doe", Age: 42}
			e := newCrypt()
			msg, err := e.Encrypt(data)
			g.Assert(err).Eql(nil)
			splitMsg := strings.Split(msg, "--")
			g.Assert(len(splitMsg)).Eql(2)
			//encryptedMsg, iv := splitMsg[0], splitMsg[1]
			var decryptedMsg Person
			err = e.Decrypt(msg, &decryptedMsg)
			g.Assert(err).Eql(nil)
			g.Assert(decryptedMsg).Eql(data)
		})

		g.It("can round trip signed and encoded string", func() {
			testData := "this is a test"
			var e MessageEncryptor
			for i := 0; i < 100; i++ {
				e = newCrypt()
				msg, err := e.EncryptAndSign(testData)
				g.Assert(err).Eql(nil)
				var output string
				err = e.DecryptAndVerify(msg, &output)
				g.Assert(err).Eql(nil)
				if output != testData {
					println(i, err.Error(), "FAILED", output, msg)
					fmt.Printf("%#v\n", e)
				}
				g.Assert(output).Eql(testData)
			}
		})

		g.It("can round trip signed and encoded struct", func() {
			e := newCrypt()
			testData := testStruct{Foo: "this is foo", Bar: 42}
			msg, err := e.EncryptAndSign(testData)
			g.Assert(err).Eql(nil)
			var output testStruct
			err = e.DecryptAndVerify(msg, &output)
			g.Assert(err).Eql(nil)
			g.Assert(output).Eql(testData)
		})
	})
}

func TestDecryptingRailsSession(t *testing.T) {
	g := Goblin(t)

	g.Describe("A Rails JSON session", func() {
		cookieContent := "TDZIdC9GcEVRSnR0aFlqYTI1SmRWTmw3NWxpRkJZNDVMK0NIUXFlcThWWitLeVQzMFVBUTE2RU82RnRsUUxQWnhyWG95dFJSRDc0OVpkVzhGWXlIb1hERHhPdk5mYStkd3pVVUZNbE1vcDRqU01MYVZJMVpMWVI5SmIweFo1N2tqWTdZcVhyWmdnc2NhZUY2b1BBMlNKWkVsT0Y0aEVQcVVKaGRISk0zR3JLWXdjaFMxamN2aThVL2hBMHBmSGx5bGg4UjUzRFErejlQVEM0eUZjcStSM3VYUkNERjBMdUVqQzZaQk5ZNHpjRT0tLUhDQ2RraWpKRDBleUp1Rm1OeVA5Snc9PQ==--61cd94a037a0a006a01403952a652ddc5da1a597"
		railsSecret := "f7b5763636f4c1f3ff4bd444eacccca295d87b990cc104124017ad70550edcfd22b8e89465338254e0b608592a9aac29025440bfd9ce53579835ba06a86f85f9"
		encryptedCookieSalt := []byte("encrypted cookie")
		encryptedSignedCookieSalt := []byte("signed encrypted cookie")

		kg := KeyGenerator{Secret: railsSecret}
		secret := kg.CacheGenerate(encryptedCookieSalt, 32)
		signSecret := kg.CacheGenerate(encryptedSignedCookieSalt, 64)
		e := MessageEncryptor{Key: secret, SignKey: signSecret}

		g.It("can be decrypted", func() {
			var session map[string]interface{}
			err := e.DecryptAndVerify(cookieContent, &session)
			g.Assert(err).Eql(nil)
			g.Assert(session["session_id"]).Eql("b2d63c07ea7a9d58e415e3672e3f31a2")
		})
	})
}

func ExampleMessageEncryptor_EncryptAndSign() {
	type Person struct {
		Id        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Age       int    `json:"age"`
	}
	john := Person{Id: 12, FirstName: "John", LastName: "Doe", Age: 42}

	k := GenerateRandomKey(32)
	signKey := []byte("this is a secret!")
	e := MessageEncryptor{Key: k, SignKey: signKey}

	// string encoding example
	msg, err := e.EncryptAndSign("my secret data")
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)

	// struct encoding example
	msg, err = e.EncryptAndSign(john)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
}

func ExampleMessageEncryptor_DecryptAndVerify() {

	type Person struct {
		Id        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Age       int    `json:"age"`
	}
	john := Person{Id: 12, FirstName: "John", LastName: "Doe", Age: 42}

	railsSecret := "f7b5763636f4c1f3ff4bd444eacccca295d87b990cc104124017ad70550edcfd22b8e89465338254e0b608592a9aac29025440bfd9ce53579835ba06a86f85f9"
	encryptedCookieSalt := []byte("encrypted cookie")
	encryptedSignedCookieSalt := []byte("signed encrypted cookie")

	kg := KeyGenerator{Secret: railsSecret}
	// use 64 bit keys since the encryption uses 32 bytes
	// but the signature uses 64. The crypto package handles that well.
	secret := kg.CacheGenerate(encryptedCookieSalt, 32)
	signSecret := kg.CacheGenerate(encryptedSignedCookieSalt, 64)
	e := MessageEncryptor{Key: secret, SignKey: signSecret}
	sessionString, err := e.EncryptAndSign(john)
	if err != nil {
		panic(err)
	}

	// decrypting the person object contained in the session
	var sessionContent Person
	err = e.DecryptAndVerify(sessionString, &sessionContent)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", sessionContent)

	//Output:
	// crypto.Person{Id:12, FirstName:"John", LastName:"Doe", Age:42}
}
