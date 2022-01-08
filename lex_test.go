package main

import (
	"fmt"
	"reflect"
	"testing"
)

type b = []byte
type ts = []token

var eq = reflect.DeepEqual

func (t token) String() string {
	return fmt.Sprintf("{%d %q}", t.typ, t.val)
}

func (toks tokenStream) flatten() (res []token) {
	for tok := range toks {
		res = append(res, tok)
	}
	return
}

func TestLex(t *testing.T) {
	tab := []struct {
		ndigits int
		data    string
		want    []token
	}{
		{0, "", nil},
		{0, "1", ts{{tokenDigits, b("1")}}},
		{1, "1", ts{{tokenDigits, b("1")}}},
		{2, "1", ts{{tokenAny, b("1")}}},
		{2, "aaa", ts{{tokenAny, b("aaa")}}},
		{3, "12", ts{{tokenAny, b("12")}}},
		{3, "123", ts{{tokenDigits, b("123")}}},
		{3, "1234", ts{{tokenDigits, b("1234")}}},
		{3, "\033[34;40m1234\033[0m", ts{
			{tokenAny, b("\033[34;40m")},
			{tokenDigits, b("1234")},
			{tokenAny, b("\033[0m")},
		}},
		{3, "\033[48;5;17m\033[38;5;19m1234\033[0m", ts{
			{tokenAny, b("\033[48;5;17m")},
			{tokenAny, b("\033[38;5;19m")},
			{tokenDigits, b("1234")},
			{tokenAny, b("\033[0m")},
		}},
		{3, "\0330000", ts{
			{tokenAny, b("\0330000")},
		}},
		{3, "\033[00x0000", ts{
			{tokenAny, b("\033[00x0000")},
		}},
	}

	for i, tc := range tab {
		have := lexTokens(b(tc.data), tc.ndigits).flatten()
		if !eq(have, tc.want) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, have, tc.want)
		}
	}
}
