package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

type Config struct {
	ndigits int
	outsep  []byte
	help    bool
}

func parseArgs(args []string) (*Config, error) {
	var (
		outsepb string
		flag    = pflag.NewFlagSet("flags", pflag.ContinueOnError)
		conf    = new(Config)
	)
	flag.SortFlags = false

	flag.StringVarP(&outsepb, "outsep", "o", " ", "output separator")
	flag.BoolVarP(&conf.help, "help", "h", false, "show this help and exit")

	flag.Usage = func() {

		fmt.Fprintln(flag.Output(),
			"Pretty-print chains of digits with a separator.")
		fmt.Fprintln(flag.Output(), "Usage: pdigit [FLAGS] N")
		fmt.Fprintln(flag.Output(), "   N\t\t\tnumber of digits to split")
		flag.PrintDefaults()
	}

	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}
	if conf.help {
		flag.SetOutput(os.Stdout)
		flag.Usage()
		return conf, nil
	}

	if len(outsepb) != 1 {
		return nil, fmt.Errorf("output separator needs to be 1 character")
	}

	conf.outsep = []byte(outsepb)

	rest := flag.Args()

	if len(rest) != 1 {
		return nil, fmt.Errorf("missing number of digits")
	}
	conf.ndigits, err = strconv.Atoi(rest[0])
	switch {
	case err != nil:
		return nil, fmt.Errorf("number of digits: %v", err)
	case conf.ndigits <= 0:
		return nil,
			fmt.Errorf("number of digist needs to be a positive number")
	}

	return conf, nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		fatal2(err)
	}
	if conf.help {
		os.Exit(0)
	}

	err = processor{conf}.run(os.Stdin, os.Stdout)
	if err != nil {
		fatal(err)
	}
}

func fatal2(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
