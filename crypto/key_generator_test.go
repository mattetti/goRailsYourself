package crypto

import (
	"fmt"
	. "github.com/franela/goblin"
	"testing"
)

func TestKegenerator_Generate(t *testing.T) {
	g := Goblin(t)
	g.Describe("Generating a derived key from a secret", func() {
		gen := KeyGenerator{Secret: "f7b5763636f4c1f3ff4bd444eacccca295d87b990cc104124017ad70550edcfd22b8e89465338254e0b608592a9aac29025440bfd9ce53579835ba06a86f85f9"}
		g.It("always generates the same key for the same salt", func() {
			salt := []byte("encrypted cookie")
			keys := make([][]byte, 10)
			for i := 0; i < 10; i++ {
				keys[i] = gen.Generate(salt, 64)
			}
			first := keys[0]
			for i := 1; i < 10; i++ {
				g.Assert(keys[i]).Eql(first)
			}
		})
	})

	g.Describe("Cache Generate a key", func() {
		gen := KeyGenerator{Secret: "f7b5763636f4c1f3ff4bd444eacccca295d87b990cc104124017ad70550edcfd22b8e89465338254e0b608592a9aac29025440bfd9ce53579835ba06a86f85f9"}
		g.It("caches the keys", func() {
			key := func(s []byte) string {
				return fmt.Sprintf("%s%d", s, 64)
			}
			salt1 := []byte("encrypted cookie")
			salt2 := []byte("signed cookie")
			_ = gen.CacheGenerate(salt1, 64)
			g.Assert(gen.cache[key(salt1)] != nil).IsTrue()
			g.Assert(gen.cache[key(salt2)] == nil).IsTrue()
			_ = gen.CacheGenerate(salt2, 64)
			g.Assert(gen.cache[key(salt2)] != nil).IsTrue()
		})
	})

}
