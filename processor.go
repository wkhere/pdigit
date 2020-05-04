package main

import (
	"bufio"
	"io"
)

type processor struct {
	*Config
}

func (p processor) transformLine(w io.Writer, input []byte) {

	for token := range lexTokens(input, p.ndigits+1) {
		switch token.typ {
		case tokenDigits:
			var i int
			l := len(token.val)

			if m := l % p.ndigits; m > 0 {
				w.Write(token.val[:m])
				w.Write(p.outsep)
				i = m
			}
			for {
				w.Write(token.val[i : i+p.ndigits])
				i += p.ndigits
				if i < l {
					w.Write(p.outsep)
				} else {
					break
				}
			}

		case tokenAny:
			w.Write(token.val)

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
