package inflector

import (
	"fmt"
	. "github.com/franela/goblin"
	"testing"
)

func ExampleParameterize() {
	fmt.Println(Parameterize("Matt Aïmonetti", "-"))
	fmt.Println(Parameterize("Ærøskøbing!", "-"))
	fmt.Println(Parameterize("Random text with *(bad)* characters", "-"))
	// Output: matt-aimonetti
	// aeroskobing
	// random-text-with-bad-characters
}

func TestParameterize(t *testing.T) {
	g := Goblin(t)
	g.Describe("Parameterize", func() {

		g.It("Should convert to lower case ", func() {
			g.Assert(Parameterize("Matt", "-")).Equal("matt")
		})

		g.It("Should transliterate", func() {
			expectations := map[string]string{
				"Ærøskøbing": "aeroskobing",
				"mon école":  "mon-ecole",
				"et ta sœur": "et-ta-soeur",
			}
			for input, output := range expectations {
				g.Assert(Parameterize(input, "-")).Equal(output)
			}
		})

		g.It("Should replace any unwanted chars by the separator", func() {
			expectations := map[string]string{
				"Matt Aimonetti":                      "matt-aimonetti",
				"Donald E. Knuth":                     "donald-e-knuth",
				"mon école":                           "mon-ecole",
				"Random text with *(bad)* characters": "random-text-with-bad-characters",
				"Trailing bad characters!@#":          "trailing-bad-characters",
				"!@#Leading bad characters":           "leading-bad-characters",
			}
			for input, output := range expectations {
				g.Assert(Parameterize(input, "-")).Equal(output)
			}
		})

		g.It("Should allow for underscore", func() {
			g.Assert(Parameterize("Allow_Under_Scores", "-")).Equal("allow_under_scores")
		})

		g.It("Should squeeze separators", func() {
			g.Assert(Parameterize("Squeeze   separators", "-")).Equal("squeeze-separators")
		})
	})
}

func ExampleTransliterate() {
	fmt.Println(Transliterate("Ærøskøbing"))
	fmt.Println(Transliterate("Ma sœur va à l'école"))
	// Output: AEroskobing
	// Ma soeur va a l'ecole
}
