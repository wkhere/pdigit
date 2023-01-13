package pdigit

import (
	"bufio"
	"io"
)

type Processor struct {
	NDigits int
	OutSep  []byte
}

func (p Processor) Run(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		p.transformLine(w, sc.Bytes())
		w.Write(LF)
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
	var i int
	l := len(digits)

	if p.NDigits <= 0 || l <= p.NDigits {
		// 1. Doesn't make sense to have ndigits zero or negative.
		//    Such setting is ignored and data is written as is.
		//    it catches the case of zero division on "empty" processor.
		// 2. If the length <= ndigits, no need to chunk; also, algo
		//    below would misbehave.
		w.Write(digits)
		return
	}

	if m := l % p.NDigits; m > 0 {
		w.Write(digits[:m])
		w.Write(p.OutSep)
		i = m
	}
	for {
		w.Write(digits[i : i+p.NDigits])
		i += p.NDigits
		if i < l {
			w.Write(p.OutSep)
		} else {
			break
		}
	}
}

var SP = []byte{0x20}
var LF = []byte{0x0a}
