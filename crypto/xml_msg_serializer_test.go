package crypto

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestXmlSerializerSerializer(t *testing.T) {
	g := Goblin(t)
	serializer := XMLMsgSerializer{}

	g.Describe("a xml serialized string", func() {
		data := "this is a test"
		output, err := serializer.Serialize(data)
		g.Assert(err).Eql(err)

		g.It("can be deserialized", func() {
			var o string
			err := serializer.Unserialize(output, &o)
			g.Assert(err).Eql(nil)
			g.Assert(o).Eql(data)
		})
	})

	g.Describe("a xml serialized struct", func() {
		type Person struct {
			Id        int    `xml:"id,attr"`
			FirstName string `xml:"name>first"`
			LastName  string `xml:"name>last"`
			Age       int    `xml:"age"`
		}
		data := Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
		output, err := serializer.Serialize(data)
		g.Assert(err).Eql(err)

		g.It("can be deserialized", func() {
			var o Person
			err := serializer.Unserialize(output, &o)
			g.Assert(err).Eql(nil)
			g.Assert(o).Eql(data)
		})
	})

}
