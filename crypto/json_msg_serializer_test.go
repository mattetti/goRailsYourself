package crypto

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestJsonMsgSerializerSerializer(t *testing.T) {
	g := Goblin(t)
	serializer := JsonMsgSerializer{}

	g.Describe("a json serialized string", func() {
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

	g.Describe("a json serialized struct", func() {
		type Person struct {
			Id        int    `json:"id"`
			FirstName string `json:"name>first"`
			LastName  string `json:"name>last"`
			Age       int    `json:"age"`
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
