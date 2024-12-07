package pdigit

import (
	"bufio"
	"fmt"
	"io"
)

const (
	defaultLineSize = 4096
	maxLineSize     = 65536
)

type reader struct {
	r      *bufio.Reader
	lineno int
}

func newReader(r io.Reader) *reader {
	return &reader{r: bufio.NewReaderSize(r, defaultLineSize)}
}

func (r *reader) ReadLine() (buf []byte, err error) {
	for {
		var b []byte

		b, err = r.r.ReadSlice('\n')
		buf = append(buf, b...)
		if err == bufio.ErrBufferFull && len(buf) < maxLineSize {
			err = nil
			continue
		}
		if len(b) > 0 {
			r.lineno++
		}
		if err != nil && err != io.EOF {
			err = fmt.Errorf("line %d: %w", r.lineno, err)
		}
		return
	}
}
