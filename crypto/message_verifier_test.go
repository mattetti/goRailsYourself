package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	. "github.com/franela/goblin"
	"strings"
	"testing"
)

type testStruct struct {
	Foo string
	Bar int
	Baz []string `json:",omitempty"`
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestMessageVerifier(t *testing.T) {
	g := Goblin(t)

	g.Describe("a malformed MessageVerifier", func() {
		g.Describe("without a serializer", func() {
			v := MessageVerifier{
				secret: []byte("Hey, I'm a secret!"),
				hasher: sha1.New,
			}

			g.It("won't generate messages", func() {
				foo := "foo"
				str, err := v.Generate(foo)
				g.Assert(err.Error()).Eql("Serializer not set")
				g.Assert(str).Eql("")
			})

			g.It("won't verify messages", func() {
				var foo string
				err := v.Verify("foo", foo)
				g.Assert(err.Error()).Eql("Serializer not set")
			})

		})

		g.Describe("without a hasher", func() {
			v := MessageVerifier{
				secret:     []byte("Hey, I'm a secret!"),
				serializer: JsonMsgSerializer{},
			}

			g.It("won't generate messages", func() {
				foo := "foo"
				str, err := v.Generate(foo)
				g.Assert(err.Error()).Eql("Hasher not set")
				g.Assert(str).Eql("")
			})

			g.It("won't verify messages", func() {
				var foo string
				err := v.Verify("foo", foo)
				g.Assert(err.Error()).Eql("Hasher not set")
			})

		})

		g.Describe("without a secret", func() {
			v := MessageVerifier{
				serializer: JsonMsgSerializer{},
				hasher:     sha1.New,
			}

			g.It("won't generate messages", func() {
				foo := "foo"
				_, err := v.Generate(foo)
				g.Assert(err.Error()).Eql("Secret not set")
			})

			g.It("won't verify messages", func() {
				var foo string
				err := v.Verify("foo", foo)
				g.Assert(err.Error()).Eql("Secret not set")
			})

		})

	})

	g.Describe("MessageVerifier with a secret & the json serializer", func() {

		g.Describe("and using SHA1", func() {
			v := MessageVerifier{
				secret:     []byte("Hey, I'm a secret!"),
				hasher:     sha1.New,
				serializer: JsonMsgSerializer{},
			}

			g.It("properly digests a string", func() {
				digest := v.DigestFor("eyJGb28iOiJmb28iLCJCYXIiOjQyfQ==")
				g.Assert(digest).Eql("b1bdb9d2b372f19dcca800e5989ee7502f1b72a5")
			})

			g.It("can do a round trip verification", func() {
				data := testStruct{Foo: "foo", Bar: 42}
				generated, err := v.Generate(data)
				g.Assert(err == nil).IsTrue()
				g.Assert(generated).Eql("eyJGb28iOiJmb28iLCJCYXIiOjQyfQ==--b1bdb9d2b372f19dcca800e5989ee7502f1b72a5")
				var verified testStruct
				err = v.Verify(generated, &verified)
				g.Assert(err == nil).IsTrue()
				g.Assert(verified).Eql(data)
			})

			g.It("can catch tampered data", func() {
				data := testStruct{Foo: "foo", Bar: 42}
				msg, err := v.Generate(data)
				// split
				dh := strings.Split(msg, "--")
				d, h := dh[0], dh[1]
				var verified testStruct
				err = v.Verify((d + "--" + h), &verified)
				g.Assert(err).Eql(nil)
				str := reverse(d) + "--" + h
				err = v.Verify(str, &verified)
				g.Assert(err.Error()).Eql("Invalid signature - bad data")
				str = d + "--" + reverse(h)
				err = v.Verify(str, &verified)
				g.Assert(err.Error()).Eql("Invalid signature - bad data")
				err = v.Verify("gargabe data", &verified)
				g.Assert(err.Error()).Eql("Invalid signature - bad data")
			})
		})

		g.Describe("and using SHA256", func() {
			v := MessageVerifier{
				secret:     []byte("Hey, I'm a secret!"),
				hasher:     sha256.New,
				serializer: JsonMsgSerializer{},
			}

			g.It("can do a round trip verification", func() {
				data := testStruct{Foo: "foo", Bar: 42}
				generated, err := v.Generate(data)
				g.Assert(err == nil).IsTrue()
				var verified testStruct
				err = v.Verify(generated, &verified)
				g.Assert(err == nil).IsTrue()
				g.Assert(verified).Eql(data)
			})
		})

		g.Describe("and using SHA512", func() {
			v := MessageVerifier{
				secret:     []byte("Hey, I'm a secret!"),
				hasher:     sha512.New,
				serializer: JsonMsgSerializer{},
			}

			g.It("can do a round trip verification", func() {
				data := testStruct{Foo: "foo", Bar: 42}
				generated, err := v.Generate(data)
				g.Assert(err == nil).IsTrue()
				var verified testStruct
				err = v.Verify(generated, &verified)
				g.Assert(err == nil).IsTrue()
				g.Assert(verified).Eql(data)
			})
		})

		g.Describe("and using md5", func() {
			v := MessageVerifier{
				secret:     []byte("Hey, I'm a secret!"),
				hasher:     md5.New,
				serializer: JsonMsgSerializer{},
			}

			g.It("can do a round trip verification", func() {
				data := testStruct{Foo: "foo", Bar: 42}
				generated, err := v.Generate(data)
				g.Assert(err == nil).IsTrue()
				var verified testStruct
				err = v.Verify(generated, &verified)
				g.Assert(err == nil).IsTrue()
				g.Assert(verified).Eql(data)
			})
		})

	})

	g.Describe("A MessageVerifier with a secret and a XML serializer", func() {

		v := MessageVerifier{
			secret:     []byte("Hey, I'm another secret!"),
			hasher:     sha1.New,
			serializer: XMLMsgSerializer{},
		}

		g.It("can do a round trip verification using SHA1", func() {
			data := testStruct{Foo: "foo", Bar: 42}
			generated, err := v.Generate(data)
			g.Assert(err == nil).IsTrue()
			var verified testStruct
			err = v.Verify(generated, &verified)
			g.Assert(err == nil).IsTrue()
			g.Assert(verified).Eql(data)
		})

	})
}

func ExampleMessageVerifier_Generate() {
	v := MessageVerifier{
		secret:     []byte("Hey, I'm a secret!"),
		hasher:     sha1.New,
		serializer: JsonMsgSerializer{},
	}
	foo := map[string]interface{}{"foo": "this is foo", "bar": 42, "baz": []string{"bar", "baz"}}
	generated, _ := v.Generate(foo)
	fmt.Println(generated)
	// Output:
	// eyJiYXIiOjQyLCJiYXoiOlsiYmFyIiwiYmF6Il0sImZvbyI6InRoaXMgaXMgZm9vIn0=--895bf35965ebef12451372225ff3f73428f48e90
}

func ExampleMessageVerifier_Verify() {
	v := MessageVerifier{
		secret:     []byte("Hey, I'm a secret!"),
		hasher:     sha1.New,
		serializer: JsonMsgSerializer{},
	}

	data := testStruct{Foo: "foo", Bar: 42}
	generated, _ := v.Generate(data)
	fmt.Println(generated)
	var verified testStruct
	_ = v.Verify(generated, &verified)
	fmt.Printf("%#v", verified)
	// Output:
	// eyJGb28iOiJmb28iLCJCYXIiOjQyfQ==--b1bdb9d2b372f19dcca800e5989ee7502f1b72a5
	// crypto.testStruct{Foo:"foo", Bar:42, Baz:[]string(nil)}
}
