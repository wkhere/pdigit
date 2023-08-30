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
		p("Pretty-print chains of digits with a separator.")
		p("Usage: pdigit [FLAGS] N")
		p("   N\t\t\t   number of digits to split")
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

	c.processor.OutSep = []byte(sep)

	rest := flag.Args()

	if len(rest) != 1 {
		return c, fmt.Errorf("missing number of digits")
	}
	c.processor.NDigits, err = strconv.Atoi(rest[0])
	switch {
	case err != nil:
		return c, fmt.Errorf("number of digits: %v", err)
	case c.processor.NDigits <= 0:
		return c,
			fmt.Errorf("number of digist needs to be a positive number")
	}

	return c, nil
}
