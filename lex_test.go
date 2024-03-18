package pdigit

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type b = []byte
type ts = []token

var eq = reflect.DeepEqual

func (t token) String() string {
	var s string
	switch t.typ {
	case tokenDigits:
		s = "tokenDigits"
	case tokenAny:
		s = "tokenAny"
	default:
		s = "!!wrong token type!!"
	}
	return fmt.Sprintf("{%s %q}", s, t.val)
}

var tabLex = []struct {
	data string
	want []token
}{
	{"", ts{}},
	{"aaa", ts{{tokenAny, b("aaa")}}},

	{"123a", ts{{tokenAny, b("123a")}}},
	{"#123", ts{{tokenAny, b("#123")}}},
	{"#a123", ts{{tokenAny, b("#a123")}}},
	{"a123", ts{
		{tokenAny, b("a")},
		{tokenDigits, b("123")},
	}},
	// ^^ todo: this should be enabled by a param, otherwise tokenAny should
	// span over the whole "a123"
	{"a#a123", ts{{tokenAny, b("a")}, {tokenAny, b("#a123")}}},

	{"1", ts{{tokenDigits, b("1")}}},
	{"12", ts{{tokenDigits, b("12")}}},
	{"123", ts{{tokenDigits, b("123")}}},
	{"1234", ts{{tokenDigits, b("1234")}}},

	{"\033[34;40m1234\033[0m", ts{
		{tokenAny, b("\033[34;40m")},
		{tokenDigits, b("1234")},
		{tokenAny, b("\033[0m")},
	}},
	{"\033[48;5;17m\033[38;5;19m1234\033[0m", ts{
		{tokenAny, b("\033[48;5;17m")},
		{tokenAny, b("\033[38;5;19m")},
		{tokenDigits, b("1234")},
		{tokenAny, b("\033[0m")},
	}},
	{"\0330000", ts{
		{tokenAny, b("\0330000")},
	}},
	{"\033[00x0000", ts{
		{tokenAny, b("\033[00x0000")},
	}},

	{"\x00", ts{{tokenAny, b("\x00")}}},
	{"\x00rest", ts{{tokenAny, b("\x00rest")}}},
	{" \x00", ts{{tokenAny, b(" ")}, {tokenAny, b("\x00")}}},
	{"aaa\x00", ts{{tokenAny, b("aaa")}, {tokenAny, b("\x00")}}},
	{"111\x00222", ts{{tokenAny, b("111\x00222")}}},
}

func TestLex(t *testing.T) {
	for i, tc := range tabLex {
		have := lexTokens(b(tc.data))
		if !eq(have, tc.want) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, have, tc.want)
		}
	}
}

func FuzzLex(f *testing.F) {
	for _, tc := range tabLex {
		f.Add(tc.data)
	}
	f.Fuzz(func(t *testing.T, s string) {
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

func BenchmarkLex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range tabLex[6:16] {
			lexTokens([]byte(tc.data))
		}
	}
}
