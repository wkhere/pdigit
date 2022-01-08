package main

import (
	"bufio"
	"io"
)

type processor struct {
	*Config
}

func (p processor) transformLine(w io.Writer, input []byte) {

	for token := range lexTokens(input) {
		switch token.typ {
		case tokenDigits:
			p.writeChunks(w, token.val)
		case tokenAny:
			w.Write(token.val)
		}
	}
}

func (p processor) writeChunks(w io.Writer, digits []byte) {
	var i int
	l := len(digits)

	if p.ndigits <= 0 || l <= p.ndigits {
		// 1. Doesn't make sense to have ndigits zero or negative.
		//    Such setting is ignored and data is written as is.
		//    it catches the case of zero division on "empty" processor.
		// 2. If the length <= ndigits, no need to chunk; also, algo
		//    below would misbehave.
		w.Write(digits)
		return
	}

	if m := l % p.ndigits; m > 0 {
		w.Write(digits[:m])
		w.Write(p.outsep)
		i = m
	}
	for {
		w.Write(digits[i : i+p.ndigits])
		i += p.ndigits
		if i < l {
			w.Write(p.outsep)
		} else {
			break
		}
	}
}

func (p processor) run(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		p.transformLine(w, sc.Bytes())
		w.Write(LF)
	}
	return sc.Err()
}

var SP = []byte{0x20}
var LF = []byte{0x0a}
