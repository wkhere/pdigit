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

func TestLex(t *testing.T) {
	tab := []struct {
		ndigits int
		data    string
		want    []token
	}{
		{0, "", nil},
		{0, "1", ts{{tokenDigits, b("1")}}},
		{1, "1", ts{{tokenDigits, b("1")}}},
		{2, "1", ts{{tokenNonDigits, b("1")}}},
		{2, "aaa", ts{{tokenNonDigits, b("aaa")}}},
		{3, "12", ts{{tokenNonDigits, b("12")}}},
		{3, "123", ts{{tokenDigits, b("123")}}},
		{3, "1234", ts{{tokenDigits, b("1234")}}},
	}

	for i, tc := range tab {
		have := lexTokens(b(tc.data), tc.ndigits).flatten()
		if !eq(have, tc.want) {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, have, tc.want)
		}
	}
}
