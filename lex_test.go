package pdigit

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type b = []byte
type ts = []token

func t(typ tokenType, val []byte) token { return token{typ: typ, val: val} }

var eq = reflect.DeepEqual

func (t token) String() string {
	var s string
	switch t.typ {
	case tokenError:
		s = "tokenError " + t.err.Error()
	case tokenDigits:
		s = "tokenDigits"
	case tokenAny:
		s = "tokenAny"
	default:
		s = "!!wrong token type!!"
	}
	return fmt.Sprintf("{%s %q}", s, t.val)
}

var tab = []struct {
	data string
	want []token
}{
	{"", ts{}},
	{"aaa", ts{t(tokenAny, b("aaa"))}},

	{"123a", ts{t(tokenAny, b("123a"))}},
	{"#123", ts{t(tokenAny, b("#123"))}},
	{"#a123", ts{t(tokenAny, b("#a123"))}},
	{"a123", ts{
		t(tokenAny, b("a")),
		t(tokenDigits, b("123")),
	}},
	// ^^ todo: this should be enabled by a param, otherwise tokenAny should
	// span over the whole "a123"
	{"a#a123", ts{t(tokenAny, b("a")), t(tokenAny, b("#a123"))}},

	{"1", ts{t(tokenDigits, b("1"))}},
	{"12", ts{t(tokenDigits, b("12"))}},
	{"123", ts{t(tokenDigits, b("123"))}},
	{"1234", ts{t(tokenDigits, b("1234"))}},

	{"\033[34;40m1234\033[0m", ts{
		t(tokenAny, b("\033[34;40m")),
		t(tokenDigits, b("1234")),
		t(tokenAny, b("\033[0m")),
	}},
	{"\033[48;5;17m\033[38;5;19m1234\033[0m", ts{
		t(tokenAny, b("\033[48;5;17m")),
		t(tokenAny, b("\033[38;5;19m")),
		t(tokenDigits, b("1234")),
		t(tokenAny, b("\033[0m")),
	}},
	{"\0330000", ts{
		t(tokenAny, b("\0330000")),
	}},
	{"\033[00x0000", ts{
		t(tokenAny, b("\033[00x0000")),
	}},

	{"\x00", ts{{tokenError, b("\x00"), lexError("binary data")}}},
	{"\x00 rest", ts{{tokenError, b("\x00"), lexError("binary data")}}},
	{" \x00", ts{
		t(tokenAny, b(" ")),
		{tokenError, b("\x00"), lexError("binary data")},
	}},
	{"aaa\x00", ts{
		t(tokenAny, b("aaa")),
		{tokenError, b("\x00"), lexError("binary data")},
	}},
}

func TestLex(t *testing.T) {
	for i, tc := range tab {
		have := lexTokens(b(tc.data))
		if !eq(have, tc.want) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, have, tc.want)
		}
	}
}

func FuzzLex(f *testing.F) {
	for _, tc := range tab {
		if containsTokenType(tc.want, tokenError) {
			continue
		}
		f.Add(tc.data)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if strings.ContainsRune(s, 0) {
			t.Skip()
		}
		var buf strings.Builder
		for _, tok := range lexTokens(b(s)) {
			buf.Write(tok.val)
		}
		res := buf.String()
		if res != s {
			t.Errorf("mismatch\nhave %v\nwant %v", res, s)
		}
	})
}

func containsTokenType(tt []token, typ tokenType) bool {
	for _, t := range tt {
		if t.typ == typ {
			return true
		}
	}
	return false
}
