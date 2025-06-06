package main

import (
	"fmt"
	"strconv"
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
	var p pstate
	var sep = " "

flags:
	for ; len(args) > 0 && p.err == nil; args = args[1:] {
		switch arg := args[0]; {

		case p.parseStringFlag(arg, "-s", "--separator", &sep):
			if len(sep) != 1 {
				return c, fmt.Errorf("output separator needs to be 1 character")
			}

		case arg == "-h", arg == "--help":
			c.help = func() { fmt.Print(usage) }
			return c, nil

		case arg == "--":
			rest = append(rest, args[1:]...)
			break flags

		case len(arg) > 1 && arg[0] == '-':
			p.errorf("unknown flag %s", arg)

		default:
			rest = append(rest, arg)
		}
	}

	if p.err != nil {
		return c, p.err
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
