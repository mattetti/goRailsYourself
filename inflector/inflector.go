package inflector

import (
	"github.com/fiam/gounidecode/unidecode"
	"regexp"
	"strings"
)

var ParameterizeReplacementRegexp = regexp.MustCompile("(?i)[^a-z0-9-_]+")

// Replaces special characters in a string so that it may be used as part of
// a 'pretty' URL.
//
func Parameterize(str, sep string) string {
	// replace accented chars with their ascii equivalents
	str = Transliterate(str)
	// Turn unwanted chars into the separator
	strB := ParameterizeReplacementRegexp.ReplaceAllLiteral([]byte(str), []byte(sep))
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
func Transliterate(str string) string {
	return unidecode.Unidecode(str)
}
