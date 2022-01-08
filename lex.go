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
	tokenAny
)

type token struct {
	typ tokenType
	val []byte
}

type tokenStream <-chan token

func lexTokens(input []byte) tokenStream {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l.tokens
}

// engine

type lexer struct {
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

const (
	cEOF rune = -1
	cESC      = '\033'
)

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

func (l *lexer) acceptOne(c rune) bool {
	if l.readc() == c {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(pred func(rune) bool) {
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
	case unicode.IsLetter(c):
		return lexLettersNoWS
	case c == cESC:
		return lexColorSeq
	default:
		l.backup()
		return lexAny
	}
}

func lexDigits(l *lexer) stateFn {
	l.acceptRun(unicode.IsNumber)
	switch next := l.peek(); {
	case next == cEOF, next == cESC, unicode.IsSpace(next):
		l.emit(tokenDigits)
		return lexStart
	default:
		return lexAny
	}
}

func lexLettersNoWS(l *lexer) stateFn {
	l.acceptRun(unicode.IsLetter)
	l.emit(tokenAny)
	return lexStart
}

func lexColorSeq(l *lexer) stateFn {
	if l.acceptOne('[') {
		return lexColorValues
	}
	return lexAny
}

func lexColorEnd(l *lexer) stateFn {
	l.emit(tokenAny)
	return lexStart
}

func lexColorValues(l *lexer) stateFn {
	l.acceptRun(unicode.IsNumber)
	switch l.readc() {
	case ';':
		return lexColorValues
	case 'm':
		return lexColorEnd
	default:
		return lexAny
	}
}

func lexAny(l *lexer) stateFn {
	l.skipUntil(unicode.IsSpace)
	l.acceptRun(unicode.IsSpace)
	l.emit(tokenAny)
	return lexStart
}
