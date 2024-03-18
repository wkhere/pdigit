package pdigit

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

var tabProc = []struct {
	spec       []int
	data, want string
}{
	{s{}, "", ""},
	{s{}, "12345", "12345"},

	{s{-1}, "", ""},
	{s{-1}, "1", "1"},
	{s{-1}, "1234", "1234"},
	{s{-1}, "12345", "12345"},

	{s{0}, "", ""},
	{s{0}, "1", "1"},
	{s{1}, "", ""},
	{s{1}, "1", "1"},
	{s{0}, "1234", "1234"},
	{s{1}, "1234", "1 2 3 4"},
	{s{2}, "1234", "12 34"},
	{s{3}, "1234", "1 234"},
	{s{4}, "1234", "1234"},
	{s{5}, "1234", "1234"},

	{s{0}, "12345", "12345"},
	{s{1}, "12345", "1 2 3 4 5"},
	{s{2}, "12345", "1 23 45"},
	{s{3}, "12345", "12 345"},
	{s{4}, "12345", "1 2345"},
	{s{5}, "12345", "12345"},
	{s{6}, "12345", "12345"},
	{s{7}, "12345", "12345"},

	{s{2}, "1", "1"},
	{s{3}, "12", "12"},

	{s{2, 4}, "", ""},
	{s{2, 4}, "1_____________", "1"},
	{s{2, 4}, "12____________", "12"},
	{s{2, 4}, "123___________", "12 3"},
	{s{2, 4}, "1234__________", "12 34"},
	{s{2, 4}, "12345_________", "12 345"},
	{s{2, 4}, "123456________", "12 3456"},
	{s{2, 4}, "1234567_______", "12 3456 7"},
	{s{2, 4}, "12345678______", "12 3456 78"},
	{s{2, 4}, "123456789_____", "12 3456 789"},
	{s{2, 4}, "1234567890____", "12 3456 7890"},
	{s{2, 4}, "12345678901___", "12 3456 7890 1"},
	{s{2, 4}, "123456789012__", "12 3456 7890 12"},
	{s{2, 4}, "1234567890123_", "12 3456 7890 123"},
	{s{2, 4}, "12345678901234", "12 3456 7890 1234"},

	{s{0, 1}, "12345678901234", "12345678901234"},
	{s{2, 0}, "12345678901234", "12 345678901234"},

	{s{-2, 1}, "12345678901234", "12345678901234"},
	{s{2, -1}, "12345678901234", "12 345678901234"},

	// todo: multiple digit tokens
}

func TestProcGood(t *testing.T) {
	for i, tc := range tabProc {
		d := strings.TrimRight(tc.data, "_")
		r := strings.NewReader(d)
		b := new(strings.Builder)

		p := Proc{GroupSpec: tc.spec, OutSep: SP}
		err := p.Run(r, b)

		if err != nil {
			t.Errorf("tc[%d] unexpected error: %v", i, err)
		}
		if have := b.String(); have != tc.want {
			t.Errorf("tc[%d] mismatch\nhave %q\nwant %q", i, have, tc.want)
		}
	}
}

type tcProcBinary struct {
	spec []int
	data string
	want string // even with error some output can be printed
}

var tabProcBinarySimple = []tcProcBinary{
	{s{2}, "\x00", ""},
	{s{2}, "123\x00", ""},
	{s{2}, "\x00123", ""},
	{s{2}, "123\x00123", ""},
	{s{2}, "123\n\x00", "1 23\n"},
}

var tabProcBinaryLong = []tcProcBinary{
	{s{2}, strings.Repeat("xx", 2047) + "\x00", ""},
	{s{2}, strings.Repeat("x ", 2047) + "\x00", ""},
	{s{2}, strings.Repeat("11", 2047) + "\x00", ""},
	{s{2}, strings.Repeat("1 ", 2047) + "\x00", ""},
}

func testProcBinary(t *testing.T, tab []tcProcBinary) {
	t.Helper()

	const errText = "binary data"

	for i, tc := range tab {
		d := strings.TrimRight(tc.data, "_")
		r := strings.NewReader(d)
		b := new(strings.Builder)

		p := Proc{GroupSpec: tc.spec, OutSep: SP}
		err := p.Run(r, b)

		if err == nil || err.Error() != errText {
			t.Errorf("tc[%d] expected error %q, have %s",
				i, errText, quote(err))
		}
		if have := b.String(); have != tc.want {
			t.Errorf("tc[%d] mismatch\nhave %q\nwant %q", i, have, tc.want)
		}
	}
}

func TestProcBinarySimple(t *testing.T) {
	testProcBinary(t, tabProcBinarySimple)
}

func TestProcBinaryLong(t *testing.T) {
	testProcBinary(t, tabProcBinaryLong)
}

func TestProcFailingWriter(t *testing.T) {
	var tab = []struct {
		spec      []int
		data      string
		failAfter int
		want      string
	}{
		{s{}, "12345", 1, "1"},

		{s{2}, "1234_", 2, "12"},
		{s{2}, "12345", 2, "1 "},

		{s{2, 4}, "1234_", 2, "12"},
		{s{2, 4}, "1234_", 3, "12 "},
		{s{2, 4}, "1234_", 4, "12 3"},

		{s{2, 4}, "12345", 2, "12"},
		{s{2, 4}, "12345", 3, "12 "},
		{s{2, 4}, "12345", 4, "12 3"},
	}

	for i, tc := range tab {
		d := strings.TrimRight(tc.data, "_")
		r := strings.NewReader(d)
		b := new(strings.Builder)
		w := &failingWriter{w: b, n: tc.failAfter}

		p := Proc{GroupSpec: tc.spec, OutSep: SP}
		err := p.Run(r, w)

		have := b.String()
		if have != tc.want {
			t.Errorf("tc#%d mismatch\nhave %q\nwant %q", i, have, tc.want)
		}
		if err == nil {
			t.Errorf("tc#%d wanted xw.err, got nil", i)
		}
	}
}

func BenchmarkProcGood(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range tabProc[15:25] {
			d := strings.TrimRight(tc.data, "_")
			r := strings.NewReader(d)
			p := Proc{GroupSpec: tc.spec, OutSep: SP}
			p.Run(r, io.Discard)
		}
	}
}

func benchProcBinary(b *testing.B, tab []tcProcBinary) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		for _, tc := range tab {
			d := strings.TrimRight(tc.data, "_")
			r := strings.NewReader(d)
			p := Proc{GroupSpec: tc.spec, OutSep: SP}
			p.Run(r, io.Discard)
		}
	}
}

func BenchmarkProcBinarySimple(b *testing.B) {
	benchProcBinary(b, tabProcBinarySimple)
}

func BenchmarkProcBinaryLong(b *testing.B) {
	benchProcBinary(b, tabProcBinaryLong)
}

func quote(x any) string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", x)
}
