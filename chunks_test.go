package pdigit

import (
	"strings"
	"testing"
)

func TestWriteChunks(t *testing.T) {
	var tab = []struct {
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
	}

	for i, tc := range tab {
		data := strings.TrimRight(tc.data, "_")
		p := Processor{GroupSpec: tc.spec, OutSep: SP}
		b := new(strings.Builder)
		p.writeChunks(&xWriter{w: b}, []byte(data))
		have := b.String()
		if have != tc.want {
			t.Errorf("tc[%d] mismatch\nhave %q\nwant %q", i, have, tc.want)
		}
	}
}

func TestWriteChunksFailingWriter(t *testing.T) {
	var tab = []struct {
		spec      []int
		data      string
		failAfter int
		want      string
	}{
		{s{}, "", 0, ""},
		{s{}, "12345", 1, "12345"},

		{s{2}, "1234_", 2, "12"},
		{s{2}, "12345", 2, "1 "},

		{s{2, 4}, "1234_", 2, "12"},
		{s{2, 4}, "1234_", 3, "12 "},
		{s{2, 4}, "1234_", 4, "12 34"},

		{s{2, 4}, "12345", 2, "12"},
		{s{2, 4}, "12345", 3, "12 "},
		{s{2, 4}, "12345", 4, "12 345"},
	}

	for i, tc := range tab {
		data := strings.TrimRight(tc.data, "_")
		p := Processor{GroupSpec: tc.spec, OutSep: SP}
		b := new(strings.Builder)
		fw := &failingWriter{w: b, n: tc.failAfter}
		xw := &xWriter{w: fw}
		p.writeChunks(xw, []byte(data))
		have := b.String()
		if have != tc.want {
			t.Errorf("tc#%d mismatch\nhave %q\nwant %q", i, have, tc.want)
		}
		if xw.err == nil {
			t.Errorf("tc#%d wanted xw.err, got nil", i)
		}
	}
}
