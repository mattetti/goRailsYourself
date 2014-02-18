// The following type alias and its functions represent
// a best effort to implement all of Ruby's String class
// functionality. As a general rule, we attempt to let
// the compiler protect us from what would normally be a
// runtime error in Ruby. And we try to return sensible defaults
// in the case where a runtime error is otherwise unavoidable.
// We also deviate from the Ruby implementation in cases where
// Go's type system is less flexible. For example:
// (golang)String#MatchIndex implements (Ruby)String#=~ but
// differs in that it returns -1 instead of nil for no matches.
// Ruby also has a much wider combination of characters to use
// for method names, as well as a more flexible system for
// passing arguments, so in some cases we remap methods to a new
// name (eg: #% maps to two methods, #Fmt and #FmtSub)
//
// Ruby and Go treat characters differently, so for our
// purposes, we will use the rune type any place the ruby docs
// expect a character.
//
// Disclaimer: This type is here for *convenience* and for when
// performance and scalability are second to readability and
// clarity of purpose.
package string

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// var r1 String
// r1 = "unamed type"
// s := "named type"
// r2 := String(s)
type String string

// Treats s as a go format specification string according to
// the docs at http://golang.org/pkg/fmt/#hdr-Printing
// @see http://ruby-doc.org/core-2.0/String.html#method-i-25
func (s String) Fmt(args ...interface{}) String {
	return String(fmt.Sprintf(string(s), args...))
}

var (
	fmtSubRegexp = regexp.MustCompile("%\\{([^\\{\\}]*)\\}")
)

// Treats s as a template and replaces each named key with
// its corresponding value. Eg: String("%{foo}").FmtNamed
// @see http://ruby-doc.org/core-2.0/String.html#method-i-25
func (s String) FmtSub(substitutions map[string]string) String {
	result := string(s)
	varSubPairs := fmtSubRegexp.FindAllStringSubmatch(string(s), -1)
	for _, varSubPair := range varSubPairs {
		subKey := varSubPair[1]
		if substitution, ok := substitutions[subKey]; ok {
			variable := varSubPair[0]
			result = strings.Replace(result, variable, substitution, -1)
		}
	}
	return String(result)
}

// Returns a new String containing n copies of the receiver.
func (s String) Repeat(n uint) String {
	result := String("")
	var i uint
	for i = 0; i < n; i++ {
		result += s
	}
	return result
}

// The given integer is first converted to it's
// code point, then appended. To concat strings or Strings
// use the + operator
// http://ruby-doc.org/core-2.0/String.html#method-i-3C-3C
func (s *String) PushI(codepoint rune) String {
	*s += String(fmt.Sprintf("%c", codepoint)) // Still easier than String("%c").Mod(t)
	return *s
}

// The given string is then appended. To concat String with
// String use the + operator
// http://ruby-doc.org/core-2.0/String.html#method-i-3C-3C
func (s *String) PushS(suffix string) String {
	*s += String(suffix)
	return *s
}

// Compares two Strings lexicographically. Returns -1, 0, or +1
// if s is less than, equal to, or greater than other.
// http://ruby-doc.org/core-2.0/String.html#method-i-3C-3D-3E
func (s String) Compare(other String) int {
	return bytes.Compare([]byte(s), []byte(other))
}

// Finds the first position a match, if any, is found. Returns -1
// if no match is found or if the pattern does not compile. This
// is a divergence from the ruby spec which returns nil in the
// case of no match.
// http://ruby-doc.org/core-2.0/String.html#method-i-3D-7E
func (s String) MatchIndex(pattern string) int {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return -1
	}
	loc := regex.FindStringIndex(string(s))
	if loc == nil {
		return -1
	} else {
		return loc[0]
	}
}

// Returns true of s matches the given pattern and
// the patterns compiles properly
// TODO: (Maybe) return a MatchData object?
// @see http://ruby-doc.org/core-2.0/String.html#method-i-match
func (s String) Matches(pattern string) bool {
	matched, _ := regexp.MatchString(pattern, string(s))
	return matched
}

// Replaces the character at the given index with substitution
// TODO: Should the arg be another String?
// @see http://ruby-doc.org/core-2.0/String.html#method-i-5B-5D-3D
func (s String) SubAt(at int, substitution string) String {
	if at < 0 || at >= len(s) {
		return s
	}
	return s[0:at] + String(substitution) + s[(at+1):]
}

// Replaces the character in s at the given index with substitution
// TODO: Should the arg be another String?
// @see http://ruby-doc.org/core-2.0/String.html#method-i-5B-5D-3D
func (s *String) SubSelfAt(at int, substitution string) {
	if at < 0 || at >= len(*s) {
		return
	}
	*s = (*s)[0:at] + String(substitution) + (*s)[(at+1):]
}

// Replaces the first match of the given pattern with substitution.
// If the pattern does not compile, then returns an unaltered copy of s.
// @see http://ruby-doc.org/core-2.0/String.html#method-i-5B-5D-3D
func (s String) Sub(pattern, substitution string) String {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	str := string(s)
	loc := regex.FindStringIndex(str)
	if loc == nil {
		return s
	}
	return String(strings.Replace(str, str[loc[0]:loc[1]], substitution, 1))
}

// Replaces in s the first match of the given pattern with substitution.
// If the pattern does not compile, then returns an unaltered copy of s.
// @see http://ruby-doc.org/core-2.0/String.html#method-i-5B-5D-3D
func (s *String) SubSelf(pattern, substitution string) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return
	}
	str := string(*s)
	loc := regex.FindStringIndex(str)
	if loc == nil {
		return
	}
	*s = String(strings.Replace(str, str[loc[0]:loc[1]], substitution, 1))
}

// Tests whether s is completely made of up valid UTF-8 characters
func (s String) IsUtf8() bool {
	return utf8.Valid([]byte(s))
}

// Tests whether s is entirely ASCII characters
// @see http://ruby-doc.org/core-2.0/String.html#method-i-ascii_only-3F
func (s String) IsASCII() bool {
	for _, b := range []byte(s) {
		if b > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-capitalize
func (s String) Capitalize() String {
	if len(s) == 0 {
		return s
	}
	return String(strings.ToUpper(string(s[0:1])) + strings.ToLower(string(s[1:])))
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-capitalize-21
func (s *String) CapitalizeSelf() {
	if len(*s) == 0 {
		return
	}
	*s = String(strings.ToUpper(string((*s)[0:1])) + strings.ToLower(string((*s)[1:])))
}

// Compare ignoring case
// @see http://ruby-doc.org/core-2.0/String.html#method-i-casecmp
func (s String) CaseCompare(other String) int {
	return s.Upcase().Compare(other.Upcase())
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-center
func (s String) Center(width int, padding ...byte) String {
	l := len(s)
	if width <= l {
		return s
	}
	leftSize := (width - l) / 2
	rightSize := leftSize + ((width - l) % 2)
	var padRstr String
	if len(padding) > 0 {
		padRstr = String(padding)
	} else {
		padRstr = " "
	}
	leftPad := padRstr.Repeat(uint(leftSize/len(padRstr)) + 1)[0:leftSize]
	rightPad := padRstr.Repeat(uint(rightSize/len(padRstr)) + 1)[0:rightSize]
	return leftPad + s + rightPad
}

// Returns string as array of chars
// @see http://ruby-doc.org/core-2.0/String.html#method-i-chars
func (s String) Chars() []String {
	chars := make([]String, len(s))
	for i, char := range strings.Split(string(s), "") {
		chars[i] = String(char)
	}
	return chars
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-chomp
func (s String) Chomp(separator ...byte) String {
	str := string(s)
	if len(separator) == 0 {
		return String(strings.TrimRight(str, "\r\n"))
	}
	return String(strings.TrimSuffix(str, string(separator)))
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-chomp-21
func (s *String) ChompSelf(separator ...byte) String {
	str := string(*s)
	if len(separator) == 0 {
		*s = String(strings.TrimRight(str, "\r\n"))
	} else {
		*s = String(strings.TrimSuffix(str, string(separator)))
	}
	return *s
}

// Chops off the trailing rune (ie: code point) or "\r\n"
// @see http://ruby-doc.org/core-2.0/String.html#method-i-chop
func (s String) Chop() String {
	l := len(s)
	if l == 0 {
		return s
	}
	if s.MatchIndex("\r\n$") >= 0 {
		return s[0 : l-2]
	}
	runes := []rune(string(s))
	return String(runes[0 : l-1])
}

// Like #chop! except that it returns the empty string
// instead of nil in the case that it is initially empty
// @see http://ruby-doc.org/core-2.0/String.html#method-i-chop-21
func (s *String) ChopSelf() String {
	l := len(*s)
	if l == 0 {
		return *s
	}
	if s.MatchIndex("\r\n$") >= 0 {
		*s = (*s)[0 : l-2]
	} else {
		*s = (*s)[0 : l-1]
	}
	return *s
}

// Returns the leading code point. If said code point spans
// multiple bytes, then the entire byte array is used.
// Hence, this is NOT a synonym for s[0:1]
// @see http://ruby-doc.org/core-2.0/String.html#method-i-chr
func (s String) Chr() String {
	l := len(s)
	if l == 0 {
		return ""
	}
	for _, c := range s {
		return String(c) // Never not
	}
	return "" // Keeping the compiler happy
}

// Empties the String
// @see http://ruby-doc.org/core-2.0/String.html#method-i-clear
func (s *String) Clear() String {
	*s = ""
	return *s
}

// Probably easier to just do []rune("my str"), but here
// for completeness' sake
// @see http://ruby-doc.org/core-2.0/String.html#method-i-codepoints
func (s String) Codepoints() []rune {
	return []rune(string(s))
}

// Each otherStr parameter defines a set of runes to count.
// The intersection of these sets defines the characters to count
// in str. Any otherStr that starts with a caret ^ is negated.
// The sequence c1-c2 means all characters between c1 and c2. The
// backslash rune can be used to escape ^ or - and is otherwise
// ignored unless it appears at the end of a sequence or the end
// of a otherStr.
//
// a = "hello world"
// a.count "lo"                   #=> 5
// a.count "lo", "o"              #=> 2
// a.count "hello", "^l"          #=> 4
// a.count "ej-m"                 #=> 4
// "hello^world".count "\\^aeiou" #=> 4
// "hello-world".count "a\\-eo"   #=> 4
// c = "hello world\\r\\n"
// c.count "\\"                   #=> 2
// c.count "\\A"                  #=> 0
// c.count "X-\\w"                #=> 3
// @see http://ruby-doc.org/core-2.0/String.html#method-i-count
// func (s String) Count(otherStr String, otherStrs ...String) int64 {
//  TODO: Implement me
// }

// @see http://ruby-doc.org/core-2.0/String.html#method-i-crypt
// func (s String) Crypt(salt String) String {
//  TODO: Implement me
// }

// @see http://ruby-doc.org/core-2.0/String.html#method-i-delete
// func (s String) Delete(otherStr String, otherStrs ...String) String {
//  TODO: Implement me
// }

// @see http://ruby-doc.org/core-2.0/String.html#method-i-delete-21
// func (s *String) DeleteSelf(otherStr String, otherStrs ...String) String {
//  TODO: Implement me
// }

// @see http://ruby-doc.org/core-2.0/String.html#method-i-downcase
func (s String) Downcase() String {
	return String(strings.ToLower(string(s)))
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-downcase-21
func (s *String) DowncaseSelf() String {
	*s = String(strings.ToLower(string(*s)))
	return *s
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-dump
func (s String) Dump() String {
	return String(strconv.Quote(string(s)))
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-each_byte
func (s String) EachByte(block func(byte)) String {
	for i := 0; i < len(s); i++ {
		block(s[i])
	}
	return s
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-each_char
func (s String) EachChar(block func(String)) String {
	for _, r := range s {
		block(String(r))
	}
	return s
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-each_codepoint
func (s String) EachCodepoint(block func(rune)) String {
	for _, r := range s {
		block(r)
	}
	return s
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-each_line
func (s String) EachLine(sep String, block func(String)) String {
	if len(s) == 0 {
		return s
	}
	if len(sep) == 0 {
		// Match paragraphs
		regx, err := regexp.Compile(`(.*\n+|.*$)`)
		if err != nil {
			return s
		}
		for _, line := range regx.FindAllString(string(s), -1) {
			block(String(line))
		}
	} else {
		for _, line := range strings.SplitAfter(string(s), string(sep)) {
			block(String(line))
		}
	}
	return s
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-empty-3F
func (s String) IsEmpty() bool {
	return len(s) == 0
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-encode
// TODO: Implement me

// @see http://ruby-doc.org/core-2.0/String.html#method-i-encode-21
// TODO: Implement me

// @see http://ruby-doc.org/core-2.0/String.html#method-i-encoding
// TODO: Implement me

// @see http://ruby-doc.org/core-2.0/String.html#method-i-end_with-3F
func (s String) EndsWith(suffix String, suffixes ...String) bool {
	if strings.HasSuffix(string(s), string(suffix)) {
		return true
	}
	for _, sfx := range suffixes {
		if strings.HasSuffix(string(s), string(sfx)) {
			return true
		}
	}
	return false
}

// @see http://ruby-doc.org/core-2.0/String.html#method-i-eql-3F
func (s String) IsEql(o String) bool {
	return s == o
}

// @see
func (s String) Upcase() String {
	return String(strings.ToUpper(string(s)))
}

var toIntRegexp = regexp.MustCompile("^[0-9]+")

func (s *String) toInt(bitSize int) int64 {
	str := toIntRegexp.FindString(string(*s))
	i, _ := strconv.ParseInt(str, 10, bitSize)
	return int64(i)
}

// Defaults to 0 if unparsable
// @see http://ruby-doc.org/core-2.0/String.html#method-i-to_i
func (s String) ToI() int {
	return int(s.toInt(0))
}

// Defaults to 0 if unparsable
// @see http://ruby-doc.org/core-2.0/String.html#method-i-to_i
func (s String) ToI8() int8 {
	return int8(s.toInt(8))
}

// Defaults to 0 if unparsable
// @see http://ruby-doc.org/core-2.0/String.html#method-i-to_i
func (s String) ToI16() int16 {
	return int16(s.toInt(16))
}

// Defaults to 0 if unparsable
// @see http://ruby-doc.org/core-2.0/String.html#method-i-to_i
func (s String) ToI32() int32 {
	return int32(s.toInt(32))
}

// Defaults to 0 if unparsable
// @see http://ruby-doc.org/core-2.0/String.html#method-i-to_i
func (s String) ToI64() int64 {
	return int64(s.toInt(64))
}

// Identical to string(s)
// @see http://ruby-doc.org/core-2.0/String.html#method-i-to_s
func (s String) ToS() string {
	return string(s)
}
