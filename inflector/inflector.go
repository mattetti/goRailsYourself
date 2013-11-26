// The inflector package ports some of Rails' ActiveSupport functions that
// can be useful outside of Rails.
//
// Rails documentation http://api.rubyonrails.org/classes/ActiveSupport/Inflector.html
package inflector

import (
	"github.com/fiam/gounidecode/unidecode"
	"regexp"
	"strings"
)

var parameterizeReplacementRegexp = regexp.MustCompile("(?i)[^a-z0-9-_]+")

// Replaces special characters in a string so that it may be used as part of
// a 'pretty' URL.
//
// Rails documentation: http://api.rubyonrails.org/classes/ActiveSupport/Inflector.html#method-i-parameterize
func Parameterize(str, sep string) string {
	// replace accented chars with their ascii equivalents
	str = Transliterate(str)
	// Turn unwanted chars into the separator
	strB := parameterizeReplacementRegexp.ReplaceAllLiteral([]byte(str), []byte(sep))
	// No more than one of the separator in a row.
	re := regexp.MustCompile(sep + `{2,}`)
	strB = re.ReplaceAllLiteral(strB, []byte(sep))
	// Remove leading/trailing separator
	re = regexp.MustCompile(`(?i)^` + sep + `|` + sep + `$`)
	strB = re.ReplaceAllLiteral(strB, []byte{})
	str = string(strB)
	// return a lower case version
	return strings.ToLower(str)
}

// Replaces non-ASCII characters with an ASCII approximation, or if none
// Transliterate("Ærøskøbing") => "AEroskobing"
// Rails documentation: http://api.rubyonrails.org/classes/ActiveSupport/Inflector.html#method-i-transliterate
func Transliterate(str string) string {
	return unidecode.Unidecode(str)
}
