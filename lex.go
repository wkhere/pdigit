package main

import (
	"unicode"
	"unicode/utf8"
)

// lexer interface

type tokenType uint

const (
	tokenError tokenType = iota
	tokenDigits
	tokenNonDigits
)

type token struct {
	typ tokenType
	val []byte
}

func lexTokens(input []byte, ndigitsAtLeast int) <-chan token {
	l := &lexer{
		ndigits: ndigitsAtLeast,
		input:   input,
		tokens:  make(chan token),
	}
	go l.run()
	return l.tokens
}

// engine

type lexer struct {
	ndigits    int
	input      []byte
	start, pos int
	lastw      int
	tokens     chan token
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for st := lexStart; st != nil; {
		st = st(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{typ: t, val: l.input[l.start:l.pos]}
	l.start = l.pos
}

// input-consuming primitives

const cEOF rune = -1

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

// input-consuming helpers

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

// state functions

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
