package main

import (
	"fmt"
	"strconv"

	"codeberg.org/wkhere/argp"
)

func parseArgs(args []string) (c config, err error) {

	const usage = `Pretty-print digits with a separator.

Usage: pdigit [FLAGS] [N1 ...] Nn

   N1 ... Nn are the groups of digits to split.

   If there is one N, print groups of N digits aligned to the right
   (12 3456 7890 for N=4).
   For N1 N2, print group of N1 then N2 digits; if there are more,
   continue N2 group while there are more digits; so effectively
   aligning from the left (12 345 678 90 for N1=2 N2=3).

Flags:
   -sC, --separator=C   output separator, one character (default " ")
   -h, --help           print this help and exit
`

	rest := make([]string, 0, len(args))
	p := argp.PState{Args: args}
	var sep = " "

flags:
	for ; len(p.Args) > 0; p.Next() {
		switch {

		case p.MatchFlagExpr("-s", "--separator"):
			if len(p.Val) != 1 {
				return c, fmt.Errorf("output separator needs to be 1 character")
			}
			sep = p.Val

		case p.MatchJustFlag("-h", "--help"):
			c.help = func() { fmt.Print(usage) }
			return c, nil

		case p.Arg() == "--":
			rest = append(rest, p.Args[1:]...)
			break flags

		case p.MatchAnyFlag():
			p.Errorf("unknown flag %s", p.Flag)

		default:
			rest = append(rest, p.Arg())
		}
	}

	if p.Err != nil {
		return c, p.Err
	}

	c.proc.OutSep = []byte(sep)

	if len(rest) < 1 {
		return c, fmt.Errorf("missing digit groups")
	}
	for _, s := range rest {
		x, err := strconv.Atoi(s)
		switch {
		case err != nil:
			return c, fmt.Errorf("error reading group spec: %w", err)
		case x <= 0:
			return c, fmt.Errorf("group spec needs to be positive: %v", x)
		}
		c.proc.GroupSpec = append(c.proc.GroupSpec, x)
	}

	return c, nil
}
