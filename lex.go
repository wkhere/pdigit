package pdigit

// lexer API

type tokenType uint

const (
	tokenDigits tokenType = iota + 1
	tokenAny
)

const tokenBufSize = 10

type token struct {
	typ tokenType
	val []byte
}

func lexTokens(input []byte) []token {
	l := &lexer{
		input:  input,
		tokens: make([]token, 0, tokenBufSize),
	}
	l.run()
	return l.tokens
}

// engine

type lexer struct {
	input      []byte
	start, pos int
	lastw      int
	tokens     []token
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for st := lexStart; st != nil; {
		st = st(l)
	}
}

func (l *lexer) emit(t tokenType) {
	l.tokens = append(l.tokens, token{typ: t, val: l.input[l.start:l.pos]})
	l.start = l.pos
}

// input-consuming primitives

const (
	cEOF rune = -1
	cESC      = '\033'
)

func (l *lexer) readc() (c rune) {
	if len(l.input[l.pos:]) == 0 {
		l.lastw = 0
		return cEOF
	}
	c = rune(l.input[l.pos])
	l.lastw = 1
	l.pos++
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
	case isDigit(c):
		return lexDigits
	case isLetter(c):
		// todo: have param to decide if there can be digits just after alpha,
		// which is now the default (for CC12 3456 7890 .. account numbers)
		return lexLetters
	case c == cESC:
		return lexColorSeq
	default:
		l.backup()
		return lexAny
	}
}

func lexDigits(l *lexer) stateFn {
	l.acceptRun(isDigit)
	switch next := l.peek(); {
	case next == cEOF, next == cESC, isSpace(next):
		l.emit(tokenDigits)
		return lexStart
	default:
		return lexAny
	}
}

func lexLetters(l *lexer) stateFn {
	l.acceptRun(isLetter)
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
	l.acceptRun(isDigit)
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
	l.skipUntil(isSpace)
	l.acceptRun(isSpace)
	l.emit(tokenAny)
	return lexStart
}

// predicates; note we consider only ascii runes

func isSpace(c rune) bool {
	switch c {
	case ' ', '\t', '\v', '\f', '\n', '\r', 0x85, 0xA0:
		return true
	}
	return false
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
