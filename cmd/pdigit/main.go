package main

import (
	"fmt"
	"io"
	"os"

	"github.com/wkhere/pdigit"
)

type config struct {
	processor pdigit.Processor

	help func(io.Writer)
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if conf.help != nil {
		conf.help(os.Stdout)
		os.Exit(0)
	}

	err = conf.processor.Run(os.Stdin, os.Stdout)
	if err != nil {
		die(1, err)
	}
}

func die(exitcode int, err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitcode)
}
