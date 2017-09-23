package main

import (
	"bufio"
	"bytes"
	"io"
)

type processor struct {
	ndigits int
}

func (p processor) transformLine(input []byte) []byte {
	var b bytes.Buffer

	for token := range lexTokens(input, p.ndigits+1) {
		switch token.typ {
		case tokenDigits:
			var i int
			l := len(token.val)

			if m := l % p.ndigits; m > 0 {
				b.Write(token.val[:m])
				b.WriteByte(' ')
				i = m
			}
			for {
				b.Write(token.val[i : i+p.ndigits])
				i += p.ndigits
				if i < l {
					b.WriteByte(' ')
				} else {
					break
				}
			}

		case tokenNonDigits:
			b.Write(token.val)

		}
	}

	return b.Bytes()
}

func (p processor) run(r io.Reader, w io.Writer) (err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		w.Write(p.transformLine(sc.Bytes()))
		w.Write(LF)
	}
	if err == io.EOF {
		err = nil
	}
	return
}

var LF = []byte{0x0a}
