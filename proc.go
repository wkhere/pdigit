package pdigit

import (
	"bufio"
	"io"
)

type Proc struct {
	GroupSpec []int
	OutSep    []byte
}

func (p Proc) Run(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)

	for {
		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && len(line) == 0 {
			break
		}

		p.transformLine(bw, line)

		if err := bw.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (p Proc) transformLine(w *bufio.Writer, input []byte) {

	for _, token := range lexTokens(input) {
		switch token.typ {
		case tokenDigits:
			p.writeChunks(w, token.val)
		case tokenAny:
			w.Write(token.val)
		}
	}
}

func (p Proc) writeChunks(w *bufio.Writer, digits []byte) {
	switch len(p.GroupSpec) {
	case 0:
		w.Write(digits)

	case 1:
		writeChunksByN(w, p.GroupSpec[0], p.OutSep, digits)

	default:
		writeChunksBySpec(w, p.GroupSpec, p.OutSep, digits)
	}
}

func writeChunksByN(w *bufio.Writer, n int, sep []byte, digits []byte) {

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
	for {
		w.Write(digits[:n])
		digits = digits[n:]
		if len(digits) == 0 {
			break
		}
		w.Write(sep)
	}
}

func writeChunksBySpec(w *bufio.Writer, spec []int, sep []byte, digits []byte) {
	// spec should have at least 2 elements, so for sure has 1
	n := spec[0]
	digits, ok := consumeDigits(w, n, digits)
	spec = spec[1:]

	for ok && len(spec) > 0 && len(digits) > 0 {
		w.Write(sep)
		n = spec[0]
		digits, ok = consumeDigits(w, n, digits)
		spec = spec[1:]
	}

	// now spec is depleted, so we use groups of last n
	for ok && len(digits) > 0 {
		w.Write(sep)
		digits, ok = consumeDigits(w, n, digits)
	}

	// similar as with writeChunksByN, if any n was zero or negative, we just
	// stop chunking now and write digits as they are
	if !ok {
		w.Write(digits)
	}
}

func consumeDigits(w *bufio.Writer, n int, digits []byte) ([]byte, bool) {
	if n <= 0 {
		return digits, false
	}
	n = min(n, len(digits))
	i, _ := w.Write(digits[:n])
	return digits[i:], true
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

var SP = []byte{0x20}
