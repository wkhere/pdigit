package pdigit

import (
	"bufio"
	"errors"
	"io"
)

type xWriter struct {
	w   io.Writer
	err error
}

func (w *xWriter) Write(p []byte) (n int, _ error) {
	if w.err != nil {
		return 0, w.err
	}
	n, w.err = w.w.Write(p)
	return n, w.err
}

type Proc struct {
	GroupSpec []int
	OutSep    []byte
}

func (p Proc) Run(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)
	xw := &xWriter{w: bw}
	for sc.Scan() {
		p.transformLine(xw, sc.Bytes())
		xw.Write(LF)
		if xw.err != nil {
			return xw.err
		}
		if err := bw.Flush(); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (p Proc) transformLine(w *xWriter, input []byte) {

	for _, token := range lexTokens(input) {
		switch token.typ {
		case tokenDigits:
			p.writeChunks(w, token.val)
		case tokenAny:
			w.Write(token.val)
		}
		if w.err != nil {
			return
		}
	}
}

func (p Proc) writeChunks(w *xWriter, digits []byte) {
	switch len(p.GroupSpec) {
	case 0:
		w.Write(digits)

	case 1:
		writeChunksByN(w, p.GroupSpec[0], p.OutSep, digits)

	default:
		writeChunksBySpec(w, p.GroupSpec, p.OutSep, digits)
	}
}

func writeChunksByN(w *xWriter, n int, sep []byte, digits []byte) {

	if n <= 0 || len(digits) <= n {
		// 1. Doesn't make sense to have n zero or negative.
		// 2. If the length <= n, no need to chunk.
		// In both cases the data is writen as it is.
		w.Write(digits)
		return
	}

	if m := len(digits) % n; m > 0 {
		w.Write(digits[:m])
		digits = digits[m:]
		w.Write(sep)
	}
	for w.err == nil {
		w.Write(digits[:n])
		digits = digits[n:]
		if len(digits) == 0 {
			break
		}
		w.Write(sep)
	}
}

func writeChunksBySpec(w *xWriter, spec []int, sep []byte, digits []byte) {
	// spec should have at least 2 elements, so for sure has 1
	n := spec[0]
	digits, err := consumeDigits(w, n, digits)
	spec = spec[1:]

	for err == nil && len(spec) > 0 && len(digits) > 0 {
		w.Write(sep)
		n = spec[0]
		digits, err = consumeDigits(w, n, digits)
		spec = spec[1:]
	}

	// now spec is depleted, so we use groups of last n
	for err == nil && len(digits) > 0 {
		w.Write(sep)
		digits, err = consumeDigits(w, n, digits)
	}

	// similar as with writeChunksByN, if any n was zero or negative, we just
	// stop chunking now and write digits as they are
	if errors.Is(err, specErr) {
		w.Write(digits)
	}
}

func consumeDigits(w *xWriter, n int, digits []byte) ([]byte, error) {
	if n <= 0 {
		return digits, specErr
	}
	n = min(n, len(digits))
	i, err := w.Write(digits[:n])
	return digits[i:], err
}

var specErr = errors.New("group spec error")

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

var SP = []byte{0x20}
var LF = []byte{0x0a}
