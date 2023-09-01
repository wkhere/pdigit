package pdigit

import (
	"bufio"
	"io"
)

type Processor struct {
	GroupSpec []int
	OutSep    []byte
}

func (p Processor) Run(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)
	for sc.Scan() {
		// todo: handle write errors
		p.transformLine(bw, sc.Bytes())
		bw.Write(LF)
		bw.Flush()
	}
	return sc.Err()
}

func (p Processor) transformLine(w io.Writer, input []byte) {

	for _, token := range lexTokens(input) {
		switch token.typ {
		case tokenDigits:
			p.writeChunks(w, token.val)
		case tokenAny:
			w.Write(token.val)
		}
	}
}

func (p Processor) writeChunks(w io.Writer, digits []byte) {
	switch len(p.GroupSpec) {
	case 0:
		w.Write(digits)

	case 1:
		writeChunksByN(w, p.GroupSpec[0], p.OutSep, digits)

	default:
		writeChunksBySpec(w, p.GroupSpec, p.OutSep, digits)
	}
}

func writeChunksByN(w io.Writer, n int, sep []byte, digits []byte) {
	var i int
	l := len(digits)

	if n <= 0 || l <= n {
		// 1. Doesn't make sense to have n zero or negative.
		// 2. If the length <= n, no need to chunk.
		// In both cases the data is writen as it is.
		w.Write(digits)
		return
	}

	if m := l % n; m > 0 {
		w.Write(digits[:m])
		w.Write(sep)
		i = m
	}
	for {
		w.Write(digits[i : i+n])
		i += n
		if i < l {
			w.Write(sep)
		} else {
			break
		}
	}
}

func writeChunksBySpec(w io.Writer, spec []int, sep []byte, digits []byte) {
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

func consumeDigits(w io.Writer, n int, digits []byte) ([]byte, bool) {
	if n <= 0 {
		return digits, false
	}
	n = min(n, len(digits))
	w.Write(digits[:n])
	return digits[n:], true
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

var SP = []byte{0x20}
var LF = []byte{0x0a}
