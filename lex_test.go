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

var tab = []struct {
	data string
	want []token
}{
	{"", ts{}},
	{"aaa", ts{{tokenAny, b("aaa")}}},

	{"123a", ts{{tokenAny, b("123a")}}},
	{"a123", ts{
		{tokenAny, b("a")},
		{tokenDigits, b("123")},
	}},
	// ^^ todo: this should be enabled by a param, otherwise tokenAny should
	// span over the whole "a123"

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
