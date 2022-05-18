package pdigit

import (
	"bytes"
	"testing"
)

func TestWriteChunks(t *testing.T) {
	var tab = []struct {
		ndigits    int
		data, want string
	}{
		{-1, "", ""},
		{-1, "1", "1"},
		{0, "", ""},
		{0, "1", "1"},
		{1, "", ""},
		{1, "1", "1"},
		{-1, "1234", "1234"},
		{0, "1234", "1234"},
		{1, "1234", "1 2 3 4"},
		{2, "1234", "12 34"},
		{3, "1234", "1 234"},
		{4, "1234", "1234"},
		{5, "1234", "1234"},
		{-1, "12345", "12345"},
		{0, "12345", "12345"},
		{1, "12345", "1 2 3 4 5"},
		{2, "12345", "1 23 45"},
		{3, "12345", "12 345"},
		{4, "12345", "1 2345"},
		{5, "12345", "12345"},
		{6, "12345", "12345"},
		{7, "12345", "12345"},
	}

	for i, tc := range tab {
		p := Processor{NDigits: tc.ndigits, OutSep: SP}
		b := new(bytes.Buffer)
		p.writeChunks(b, []byte(tc.data))
		have := b.String()
		if have != tc.want {
			t.Errorf("tc[%d] mismatch\nhave %v\nwant %v", i, have, tc.want)
		}
	}
}
