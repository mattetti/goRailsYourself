package crypto

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestNullSerializerSerializer(t *testing.T) {
	g := Goblin(t)
	serializer := NullMsgSerializer{}

	g.Describe("a null serialized string", func() {
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

	g.Describe("a null serialized struct", func() {
		data := map[string]string{"foo": "matt", "bar": "aimonetti"}
		output, err := serializer.Serialize(data)

		g.It("serializes properly", func() {
			g.Assert(err).Eql(err)
			g.Assert(output).Eql("map[bar:aimonetti foo:matt]")
		})

		g.It("can be deserialized", func() {
			var o string
			err := serializer.Unserialize(output, &o)
			g.Assert(err).Eql(nil)
			g.Assert(o).Eql("map[bar:aimonetti foo:matt]")
		})
	})

}
