package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/wkhere/pdigit"
)

type config struct {
	processor pdigit.Processor
	help    bool
}

func parseArgs(args []string) (*config, error) {
	var (
		sep  string
		flag = pflag.NewFlagSet("flags", pflag.ContinueOnError)
		conf = new(config)
	)
	flag.SortFlags = false

	flag.StringVarP(&sep, "separator", "s", " ", "output separator")
	flag.BoolVarP(&conf.help, "help", "h", false, "show this help and exit")

	flag.Usage = func() {

		fmt.Fprintln(flag.Output(),
			"Pretty-print chains of digits with a separator.")
		fmt.Fprintln(flag.Output(), "Usage: pdigit [FLAGS] N")
		fmt.Fprintln(flag.Output(), "   N\t\t\t   number of digits to split")
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

	if len(sep) != 1 {
		return nil, fmt.Errorf("output separator needs to be 1 character")
	}

	conf.processor.OutSep = []byte(sep)

	rest := flag.Args()

	if len(rest) != 1 {
		return nil, fmt.Errorf("missing number of digits")
	}
	conf.processor.NDigits, err = strconv.Atoi(rest[0])
	switch {
	case err != nil:
		return nil, fmt.Errorf("number of digits: %v", err)
	case conf.processor.NDigits <= 0:
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

	err = conf.processor.Run(os.Stdin, os.Stdout)
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
