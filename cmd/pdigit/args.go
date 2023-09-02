package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/pflag"
)

func parseArgs(args []string) (c config, err error) {
	var (
		flag = pflag.NewFlagSet("flags", pflag.ContinueOnError)
		sep  string
		help bool
	)
	flag.SortFlags = false

	flag.StringVarP(&sep, "separator", "s", " ", "output separator")
	flag.BoolVarP(&help, "help", "h", false, "show this help and exit")

	flag.Usage = func() {
		p := func(s string) { fmt.Fprintln(flag.Output(), s) }
		p("Pretty-print digits with a separator.")
		p("Usage: pdigit [FLAGS] [N1 ...] Nn")
		p("   N1 ... Nn\t\t   groups of digits to split")
		p("   If there is one N, print groups of N digits aligned to the right")
		p("   (12 3456 7890 for N=4).")
		p("   For N1 N2, print group of N1 then N2 digits; if there are more,")
		p("   continue N2 group while there are more digits; so effectively")
		p("   aligning from the left (12 345 678 90 for N1=2 N2=4).")
		p("Flags:")
		flag.PrintDefaults()
	}

	err = flag.Parse(args)
	if err != nil {
		return c, err
	}
	if help {
		c.help = func(w io.Writer) {
			flag.SetOutput(w)
			flag.Usage()
		}
		return c, nil
	}

	if len(sep) != 1 {
		return c, fmt.Errorf("output separator needs to be 1 character")
	}

	c.proc.OutSep = []byte(sep)

	rest := flag.Args()

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
