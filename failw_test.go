package pdigit

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type failingWriter struct {
	w    io.Writer
	n, j int // fails after n bytes written

}

func (w *failingWriter) Write(p []byte) (int, error) {
	i, err := w.w.Write(p)
	w.j += i
	if w.j >= w.n && err == nil {
		err = specialErr
	}
	return i, err
}

var specialErr = errors.New("the time has come")

func TestFailingWriterBasic(t *testing.T) {
	b := new(strings.Builder)
	w := &failingWriter{w: b, n: 4}

	_, err := io.WriteString(w, "foo")
	if err != nil {
		t.Errorf("want nil, have err `%v`", err)
	}
	_, err = io.WriteString(w, "1")
	if !errors.Is(err, specialErr) {
		t.Errorf("want err `%v`, have err `%v`", specialErr, err)
	}
}
