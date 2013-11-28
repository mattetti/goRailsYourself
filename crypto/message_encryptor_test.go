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
		signKey := "this is a secret!"
		e := MessageEncryptor{key: k, signKey: signKey}
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
			//g.Assert(newMsg).Eql("my secret data")
		})

	})

}

func TestMessageEncryptor(t *testing.T) {
	g := Goblin(t)

	g.Describe("MessageEncryptor properly setup using aes cbc", func() {
		newCrypt := func() MessageEncryptor {
			return MessageEncryptor{key: GenerateRandomKey(32),
				cipher: "aes-cbc",
				verifier: &MessageVerifier{
					secret:     "signature secret!",
					hasher:     sha1.New,
					serializer: NullMsgSerializer{},
				},
				serializer: JsonMsgSerializer{},
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

func ExampleMessageEncryptor_EncryptAndSign() {
	type Person struct {
		Id        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Age       int    `json:"age"`
	}
	john := Person{Id: 12, FirstName: "John", LastName: "Doe", Age: 42}

	k := GenerateRandomKey(32)
	signKey := "this is a secret!"
	e := MessageEncryptor{key: k, signKey: signKey}

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

	k := GenerateRandomKey(32)
	signKey := "this is a secret!"
	e := MessageEncryptor{key: k, signKey: signKey}
	// string encoding example
	encryptedString, err := e.EncryptAndSign("my secret data")
	if err != nil {
		panic(err)
	}
	// struct encoding example
	encryptedPerson, err := e.EncryptAndSign(john)
	if err != nil {
		panic(err)
	}

	// decrypting a string message
	var decryptedString string
	err = e.DecryptAndVerify(encryptedString, &decryptedString)
	if err != nil {
		panic(err)
	}
	fmt.Println(decryptedString)

	// decrypting the person object
	var decryptedPerson Person
	err = e.DecryptAndVerify(encryptedPerson, &decryptedPerson)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", decryptedPerson)

	//Output:
	// my secret data
	// crypto.Person{Id:12, FirstName:"John", LastName:"Doe", Age:42}
}
