package main

import (
	"fmt"
	"strings"
)

type pstate struct{ err error }

func (p *pstate) errorf(format string, a ...any) bool {
	// saving only the first error
	if p.err == nil {
		p.err = fmt.Errorf(format, a...)
	}
	return false
}

func (p *pstate) eatFlagPrefix(s, short, long string) (flag, v string, ok bool) {
	var n int
	var needsEq bool

	if strings.HasPrefix(s, long) {
		n, needsEq = len(long), true
	} else if strings.HasPrefix(s, short) {
		n = len(short)
	} else {
		return "", "", false
	}

	flag, s = s[:n], s[n:]
	if needsEq {
		if len(s) == 0 || s[0] != '=' {
			return flag, "", p.errorf("flag %s needs '=' and a value", flag)
		}
		s = s[1:]
	} else {
		if len(s) > 0 && s[0] == '=' {
			s = s[1:]
		}
	}

	if len(s) == 0 {
		return flag, "", p.errorf("flag %s needs a value", flag)
	}
	return flag, s, true
}

func (p *pstate) parseStringFlag(s, short, long string, dest *string) bool {
	_, s, ok := p.eatFlagPrefix(s, short, long)
	if !ok {
		return false
	}
	*dest = s
	return true
}
