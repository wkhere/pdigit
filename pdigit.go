package main

// todo:
// - various separators for numbers in input (not only spaces)
// - output separator for number chunks given by flag

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type tokenType uint

const (
	tokenError tokenType = iota
	tokenDigits
	tokenNonDigits
)

const cEOF rune = -1

type token struct {
	typ tokenType
	val []byte
}

type lexer struct {
	ndigits    int
	input      []byte
	start, pos int
	lastw      int
	tokens     chan token
}

type stateFn func(*lexer) stateFn

func lexTokens(input []byte, ndigitsAtLeast int) <-chan token {
	l := &lexer{
		ndigits: ndigitsAtLeast,
		input:   input,
		tokens:  make(chan token),
	}
	go l.run()
	return l.tokens
}

func (l *lexer) readc() (c rune) {
	c, l.lastw = utf8.DecodeRune(l.input[l.pos:])
	if l.lastw == 0 {
		return cEOF
	}
	l.pos += l.lastw
	return c
}

// backup can be used only once after each readc.
func (l *lexer) backup() {
	l.pos -= l.lastw
}

func (l *lexer) peek() rune {
	c := l.readc()
	l.backup()
	return c
}

func (l *lexer) acceptAny(pred func(rune) bool) {
	for pred(l.readc()) {
	}
	l.backup()
}

func (l *lexer) skipUntil(pred func(rune) bool) {
	for {
		if c := l.readc(); c == cEOF || pred(c) {
			break
		}
	}
	l.backup()
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{typ: t, val: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) run() {
	for st := lexStart; st != nil; {
		st = st(l)
	}
	close(l.tokens)
}

func lexStart(l *lexer) stateFn {
	switch c := l.readc(); {
	case c == cEOF:
		return nil
	case unicode.IsNumber(c):
		return lexDigits
	default:
		l.backup()
		return lexNonDigits
	}
}

func lexDigits(l *lexer) stateFn {
	l.acceptAny(unicode.IsNumber)
	if next := l.peek(); next == cEOF || unicode.IsSpace(next) {
		if l.pos-l.start >= l.ndigits {
			l.emit(tokenDigits)
			return lexStart
		}
	}
	return lexNonDigits
}

func lexNonDigits(l *lexer) stateFn {
	l.skipUntil(unicode.IsSpace)
	l.acceptAny(unicode.IsSpace)
	l.emit(tokenNonDigits)
	return lexStart
}

func (p processor) transformLine(input []byte) []byte {
	b := bytes.NewBuffer(nil)

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

type processor struct {
	ndigits int
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

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	ndigits, err := strconv.Atoi(os.Args[1])
	if err != nil || ndigits <= 0 {
		usage()
	}

	err = processor{ndigits}.run(os.Stdin, os.Stdout)
	if err != nil {
		fatal(err)
	}
}

var LF = []byte{0x0a}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: pdigit n")
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
