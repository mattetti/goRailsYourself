package string

import (
	"fmt"
	"math"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestString(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe(".String is the type initializer", func() {
		g.It("should be initializable from string", func() {
			str := "named type"
			s := String(str)
			Expect(fmt.Sprint(s)).To(Equal(fmt.Sprint(str)))
		})
		g.It("should be assignable from a string", func() {
			var str String
			str = "assigned"
			Expect(str).To(Equal(String("assigned")))
		})
	})

	// #% (format spec case)
	g.Describe("#Fmt is like ruby's String#%", func() {
		g.It("should format input arguments", func() {
			// TODO: Move these to examples
			zeroPadded := String("%05d").Fmt(123)
			Expect(zeroPadded).To(Equal(String("00123")))
			paddedSuffixAndHexed := String("%-5s: %08x").Fmt("ID", 537793750)
			Expect(paddedSuffixAndHexed).To(Equal(String("ID   : 200e14d6")))

			// Formatting:
			s := struct{ a string }{a: "a"}
			formatted := String("%v %+v %T %t %b %f %s").Fmt(s, s, String(" "), true, 3, 1.23, String("string_as_string"))
			Expect(formatted).To(Equal(String("{a} {a:a} string.String true 11 1.230000 string_as_string")))
		})
	})

	// #% (name replacement case)
	g.Describe("#FmtSub is like ruby's String#%", func() {
		g.It("should replace with substitutions, skipping unmatched names", func() {
			// Name replacement:
			nameReplaced := String("foo = %{foo} %{foo} %{baz} %{doge}").FmtSub(map[string]string{"foo": "bar", "baz": "biz"})
			Expect(nameReplaced).To(Equal(String("foo = bar bar biz %{doge}")))
		})
	})

	// #*
	g.Describe("#Repeat is like ruby's String#*", func() {
		g.It("should return new string with multiple copies of the receiver", func() {
			Expect(String("Ho! ").Repeat(3)).To(Equal(String("Ho! Ho! Ho! ")))
			Expect(String("Ho! ").Repeat(0)).To(Equal(String("")))
		})
	})

	g.Describe("#PushI is like ruby's String#<<(integer) or String#concat()", func() {
		g.It("should append integers as codepoints", func() {
			s := String("exciting")
			r := s.PushI(33)
			Expect(s).To(Equal(String("exciting!")))
			Expect(&s).To(Equal(&r))
		})
	})

	g.Describe("#PushS is like ruby's String#<<(string) or String#concat(", func() {
		g.It("should append integers as codepoints", func() {
			s := String("exciting")
			r := s.PushS(", isn't it?")
			Expect(s).To(Equal(String("exciting, isn't it?")))
			Expect(&s).To(Equal(&r))
		})
	})

	g.Describe("#Compare is like ruby's String#<=>", func() {
		g.It("should compare lexicographically", func() {
			Expect(String("abcdef").Compare(String("abcde"))).To(Equal(1))
			Expect(String("abcdef").Compare(String("abcdef"))).To(Equal(0))
			Expect(String("abcdef").Compare(String("abcdefg"))).To(Equal(-1))
			Expect(String("abcdef").Compare(String("ABCDEF"))).To(Equal(1))
		})
	})

	g.Describe("#MatchIndex is like ruby's String#=~", func() {
		g.It("should find the first index of the regex", func() {
			Expect(String("hi ho ho").MatchIndex("ho")).To(Equal(3))
		})
		g.It("should return -1 if no match is found", func() {
			Expect(String("hi ho ho").MatchIndex("he")).To(Equal(-1))
		})
	})

	g.Describe("#Matches is like ruby's String#match but only returns true/false", func() {
		g.It("should indicated if pattern is matched or not", func() {
			Expect(String("hi ho ho").Matches("ho")).To(Equal(true))
			Expect(String("hi ho ho").Matches("hey")).To(Equal(false))
		})
	})

	g.Describe("#SubAt is like ruby's String#[index]= but returns a copy", func() {
		g.It("should replace at the specified index and return new String", func() {
			s := String("fee i fo fum")
			s = s.SubAt(4, "fi")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should replace the last element if index is len - 1", func() {
			s := String("fee fi fo m")
			s = s.SubAt(10, "fum")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should no-op if index out of bounds", func() {
			s := String("fee fi fo fum")
			s = s.SubAt(100, "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
			s = s.SubAt(-1, "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
	})

	g.Describe("#SubSelfAt is like ruby's String#[index]=", func() {
		g.It("should modify the String in place", func() {
			s := String("fee i fo fum")
			s.SubSelfAt(4, "fi")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should modify even the last element", func() {
			s := String("fee fi fo m")
			s.SubSelfAt(10, "fum")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should no-op if index out of bounds", func() {
			s := String("fee fi fo fum")
			s.SubSelfAt(100, "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
			s.SubSelfAt(-1, "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
	})

	g.Describe("#Sub is like ruby's String#sub", func() {
		g.It("should replace at the specified pattern and return new String", func() {
			s := String("fee i fo fum")
			s = s.Sub(" i ", " fi ")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should properly compile regex patterns", func() {
			s := String("fee fi fo m")
			s = s.Sub(".$", "fum")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should no-op if no match found", func() {
			s := String("fee fi fo fum")
			s = s.Sub("foo", "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
			s = s.Sub("a(.", "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
	})

	g.Describe("#SubSelf is like ruby's String#sub!", func() {
		g.It("should replace at the specified pattern and return new String", func() {
			s := String("fee i fo fum")
			s.SubSelf(" i ", " fi ")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should properly compile regex patterns", func() {
			s := String("fee fi fo m")
			s.SubSelf(".$", "fum")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
		g.It("should no-op if no match found", func() {
			s := String("fee fi fo fum")
			s.SubSelf("foo", "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
			s.SubSelf("a(.", "bar")
			Expect(s).To(Equal(String("fee fi fo fum")))
		})
	})

	g.Describe("#IsUtf8 determines if all characters are utf8 (non-ruby method)", func() {
		g.It("should return true if all characters are utf8", func() {
			s := String("Hello, 世界")
			Expect(s.IsUtf8()).To(BeTrue())
		})
		g.It("should return false if not all characters are utf8", func() {
			s := String(string([]byte{0xff, 0xfe, 0xfd}))
			Expect(s.IsUtf8()).To(BeFalse())
		})
	})

	g.Describe("#IsASCII is like ruby's String#ascii_only?", func() {
		g.It("should return true if all characters are ASCII", func() {
			s := String("Hello, World")
			Expect(s.IsASCII()).To(BeTrue())
		})
		g.It("should return false if not all characters are ASCII", func() {
			s := String("Hello, 世界")
			Expect(s.IsASCII()).To(BeFalse())
		})
	})

	g.Describe("#Capitalize is like ruby's String#capitalize", func() {
		g.It("should uppercase the first letter and lower case the rest", func() {
			Expect(String("hello").Capitalize()).To(Equal(String("Hello")))
			Expect(String("HELLO").Capitalize()).To(Equal(String("Hello")))
			Expect(String("123ABC").Capitalize()).To(Equal(String("123abc")))
			Expect(String("").Capitalize()).To(Equal(String("")))
		})
	})

	g.Describe("#CapitalizeSelf is like ruby's String#capitalize!", func() {
		g.It("should uppercase the first letter and lower case the rest", func() {
			s := String("hello")
			s.CapitalizeSelf()
			Expect(s).To(Equal(String("Hello")))
			s = String("HELLO")
			s.CapitalizeSelf()
			Expect(s).To(Equal(String("Hello")))
			s = String("123ABC")
			s.CapitalizeSelf()
			Expect(s).To(Equal(String("123abc")))
			s = String("")
			s.CapitalizeSelf()
			Expect(s).To(Equal(String("")))
		})
	})

	g.Describe("#CaseCompare is like ruby's String#casecmp(string)", func() {
		g.It("should compare lexicographically ignoring case", func() {
			left := String("aBcDef")
			right := String("AbCdE")
			Expect(left.CaseCompare(right)).To(Equal(1))
			left = String("aBcDef")
			right = String("AbCdEf")
			Expect(left.CaseCompare(right)).To(Equal(0))
			left = String("aBcDef")
			right = String("AbCdEfG")
			Expect(left.CaseCompare(right)).To(Equal(-1))
			left = String("aBcDef")
			right = String("ABCDEF")
			Expect(left.CaseCompare(right)).To(Equal(0))
		})
	})

	g.Describe("#Center is like ruby's String#center(width, padstr=nil)", func() {
		g.It("should center the text in whitespace", func() {
			s := String("center")
			Expect(s.Center(5)).To(Equal(s))
			Expect(s.Center(7)).To(Equal(String("center ")))
			Expect(s.Center(16)).To(Equal(String("     center     ")))
		})
		g.It("should center the text in padding", func() {
			s := String("center")
			Expect(s.Center(5, '-')).To(Equal(s))
			Expect(s.Center(7, "-"[0])).To(Equal(String("center-")))
			Expect(s.Center(16, []byte("123")...)).To(Equal(String("12312center12312")))
		})
	})

	g.Describe("#Chars is like ruby's String#chars", func() {
		g.It("should split String into slice of chars", func() {
			Expect(String("abc").Chars()).To(Equal([]String{"a", "b", "c"}))
		})
	})

	g.Describe("#Chomp is like ruby's String#chomp", func() {
		g.It("should chomp off trailing newlines", func() {
			Expect(String("hello \r\n").Chomp()).To(Equal(String("hello ")))
		})
		g.It("should chomp off trailing separator", func() {
			Expect(String("hello, world日").Chomp([]byte(", world日")...)).To(Equal(String("hello")))
		})
	})

	g.Describe("#ChompSelf is like ruby's String#chomp!", func() {
		g.It("should chomp off trailing newlines", func() {
			s := String("hello \r\n")
			r := s.ChompSelf()
			Expect(s).To(Equal(String("hello ")))
			Expect(&s).To(Equal(&r))
		})
		g.It("should chomp off trailing separator", func() {
			s := String("hello, world")
			r := s.ChompSelf([]byte(", world")...)
			Expect(s).To(Equal(String("hello")))
			Expect(&s).To(Equal(&r))
		})
	})

	g.Describe("#Chop is like ruby's String#chop", func() {
		g.It("should be a no-op for the empty string", func() {
			Expect(String("").Chop()).To(Equal(String("")))
		})
		g.It("should chop off \\r\\n", func() {
			Expect(String("foo\r\n").Chop()).To(Equal(String("foo")))
		})
		g.It("should chop off trailing characters", func() {
			Expect(String("foo bar").Chop()).To(Equal(String("foo ba")))
		})
	})

	g.Describe("#ChopSelf is like ruby's String#chop!", func() {
		g.It("should be a no-op for the empty string", func() {
			s := String("")
			r := s.ChopSelf()
			Expect(s).To(Equal(String("")))
			Expect(&s).To(Equal(&r))
		})
		g.It("should chop off \\r\\n", func() {
			s := String("chop these\r\n")
			r := s.ChopSelf()
			Expect(s).To(Equal(String("chop these")))
			Expect(&s).To(Equal(&r))
		})
		g.It("should chop off trailing characters", func() {
			s := String("chop it")
			r := s.ChopSelf()
			Expect(s).To(Equal(String("chop i")))
			Expect(&s).To(Equal(&r))
		})
	})

	g.Describe("#Chr", func() {
		g.It("should return empty if called on empty String", func() {
			Expect(String("").Chr()).To(Equal(String("")))
		})
		g.It("should return leading code point if called on non-empty String", func() {
			Expect(String("foo").Chr()).To(Equal(String("f")))
		})
		g.It("should return leading code point if code point is multiple bytes", func() {
			Expect(String("日").Chr()).To(Equal(String("日")))
		})
	})

	g.Describe("#Clear is like ruby's String#clear", func() {
		g.It("should empty the String", func() {
			s := String("clear me")
			r := s.Clear()
			Expect(s).To(Equal(String("")))
			Expect(&s).To(Equal(&r))
		})
	})

	g.Describe("#Codepoints is like ruby's String#codepoints", func() {
		g.It("should return an array of runes representing string", func() {
			Expect(String("<日>").Codepoints()).To(Equal([]rune{60, 26085, 62}))
		})
	})

	// TODO: Fixme
	// g.Describe("#Count is like ruby's String#count", func() {
	//   s := String("hello world")
	//   g.It("should count all runes in a regular string", func() {
	//     Expect(s.Count("lo")).To(Equal(int64(5)))
	//   })
	//   g.It("should count only the intersection of runes", func() {
	//     Expect(s.Count("lo", "o")).To(Equal(int64(2)))
	//   })
	//   g.It("should parse '-' as range of runes", func() {
	//     Expect(s.Count("ej-m")).To(Equal(int64(4)))
	//   })
	//   g.It("should interpret '^' as exclusionary", func() {
	//     Expect(s.Count("hello", "^l")).To(Equal(int64(4)))
	//   })
	//   s = String("hello^world")
	//   g.It("should count '^', ", func() {
	//     Expect(s.Count("\\^aeiou")).To(Equal(int64(4)))
	//   })
	//   g.It("should ", func() {
	//     Expect(s.Count("\\^aeiou")).To(Equal(int64(4)))
	//   })
	//   s = String("hello world\\r\\n")
	//   g.It("should ", func() {

	//   })
	//   g.It("should ", func() {

	//   })
	//   g.It("should ", func() {

	//   })
	// })

	g.Describe("#Downcase is like ruby's String#downcase", func() {
		g.It("should convert all characters to their lowercase equivalent", func() {
			Expect(String("AbCd123-$").Downcase()).To(Equal(String("abcd123-$")))
		})
	})

	g.Describe("#DowncaseSelf is like ruby's String#downcase!", func() {
		g.It("should convert all characters to their lowercase equivalent, in place", func() {
			s := String("AbCd123-$")
			r := s.DowncaseSelf()
			Expect(s).To(Equal(String("abcd123-$")))
			Expect(&s).To(Equal(&r))
		})
	})

	g.Describe("#Dump is like ruby's String#dump", func() {
		g.It("should replace all non-printing characters by \\n notation and escape special characters", func() {
			Expect(String("hello \n ''").Dump()).To(Equal(String("\"hello \\n ''\"")))
		})
	})

	g.Describe("#EachByte is like ruby's String#each_byte", func() {
		g.It("should pass each byte to the given function", func() {
			bytes := []byte{}
			block := func(b byte) { bytes = append(bytes, b) }
			Expect(String("hello 日").EachByte(block)).To(Equal(String("hello 日")))
			Expect(bytes).To(Equal([]byte{104, 101, 108, 108, 111, 32, 230, 151, 165}))
		})
	})

	g.Describe("#EachChar is like ruby's String#each_char", func() {
		g.It("should pass each byte to the given function", func() {
			chars := []String{}
			block := func(c String) { chars = append(chars, c) }
			Expect(String("hello 日").EachChar(block)).To(Equal(String("hello 日")))
			Expect(chars).To(Equal([]String{"h", "e", "l", "l", "o", " ", "日"}))
		})
	})

	g.Describe("#EachCodepoint is like ruby's String#each_codepoint", func() {
		g.It("should pass each byte to the given function", func() {
			runes := []rune{}
			block := func(r rune) { runes = append(runes, r) }
			Expect(String("hello 日").EachCodepoint(block)).To(Equal(String("hello 日")))
			Expect(runes).To(Equal([]rune{'h', 'e', 'l', 'l', 'o', ' ', '日'}))
		})
	})

	g.Describe("#EachLine is like ruby's String#each_line", func() {
		g.It("should ", func() {
			lines := []String{}
			block := func(line String) { lines = append(lines, line) }
			Expect(String("hello\nworld").EachLine("\n", block)).To(Equal(String("hello\nworld")))
			Expect(lines).To(Equal([]String{"hello\n", "world"}))
			lines = []String{}
			Expect(String("hello\nworld").EachLine("l", block)).To(Equal(String("hello\nworld")))
			Expect(lines).To(Equal([]String{"hel", "l", "o\nworl", "d"}))
			lines = []String{}
			Expect(String("hello\n\n\nworld\nguys").EachLine("", block)).To(Equal(String("hello\n\n\nworld\nguys")))
			Expect(lines).To(Equal([]String{"hello\n\n\n", "world\n", "guys"}))
		})
	})

	g.Describe("#IsEmpty is like ruby's String#empty?", func() {
		g.It("should return true if length is zero", func() {
			Expect(String("").IsEmpty()).To(BeTrue())
			Expect(String("K").IsEmpty()).To(BeFalse())
		})
	})

	g.Describe("#EndsWith is like ruby's String#ends_with?", func() {
		g.It("should be true if any suffix matches", func() {
			Expect(String("very shibe").EndsWith("wow")).To(BeFalse())
			Expect(String("very shibe").EndsWith("shibe")).To(BeTrue())
			Expect(String("very shibe").EndsWith("wow", "shibe", "much", "doge")).To(BeTrue())
		})
	})

	g.Describe("#IsEql is like ruby's String#eql?", func() {
		g.It("should test for equality", func() {
			Expect(String("money on my mind").IsEql(String("mind on my money"))).To(BeFalse())
			Expect(String("herby hancock").IsEql(String("herby hancock"))).To(BeTrue())
		})
	})

	g.Describe("#", func() {
		g.It("should ", func() {

		})
	})

	g.Describe("#ToI is like ruby's String#to_i", func() {
		g.It("should parse integer prefixes", func() {
			s := String(fmt.Sprintf("%vx", math.MaxInt8))
			Expect(s.ToI()).To(Equal(int(127)))
		})
		g.It("should fail to 0", func() {
			s := String("x12")
			Expect(s.ToI()).To(Equal(int(0)))
		})
	})

	g.Describe("#ToI8 is like ruby's String#to_i", func() {
		g.It("should parse integers prefixes and fail to 0", func() {
			s := String(fmt.Sprintf("%vx", math.MaxInt8))
			Expect(s.ToI8()).To(Equal(int8(127)))
		})
		g.It("should fail to 0", func() {
			s := String("x12")
			Expect(s.ToI8()).To(Equal(int8(0)))
		})
	})

	g.Describe("#ToI16 is like ruby's String#to_i", func() {
		g.It("should parse integers prefixes and fail to 0", func() {
			s := String(fmt.Sprintf("%vx", math.MaxInt16))
			Expect(s.ToI16()).To(Equal(int16(32767)))
		})
		g.It("should fail to 0", func() {
			s := String("x12")
			Expect(s.ToI16()).To(Equal(int16(0)))
		})
	})

	g.Describe("#ToI32 is like ruby's String#to_i", func() {
		g.It("should parse integers prefixes and fail to 0", func() {
			s := String(fmt.Sprintf("%vx", math.MaxInt32))
			Expect(s.ToI32()).To(Equal(int32(2147483647)))
		})
		g.It("should fail to 0", func() {
			s := String("x12")
			Expect(s.ToI32()).To(Equal(int32(0)))
		})
	})

	g.Describe("#ToI64 is like ruby's String#to_i", func() {
		g.It("should parse integers prefixes and fail to 0", func() {
			s := String(fmt.Sprintf("%vx", math.MaxInt64))
			Expect(s.ToI64()).To(Equal(int64(9223372036854775807)))
		})
		g.It("should fail to 0", func() {
			s := String("x12")
			Expect(s.ToI64()).To(Equal(int64(0)))
		})
	})

	g.Describe("#ToS is like ruby's String#to_s", func() {
		g.It("should be equal, char for char, to it's equivalent string", func() {
			str := "string"
			Expect([]byte(String(str))).To(Equal([]byte(str)))
		})
	})
}
