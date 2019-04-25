package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	ndigits, err := strconv.Atoi(os.Args[1])
	if err != nil || ndigits <= 0 {
		usage()
	}

	err = processor{ndigits}.run(os.Stdin, os.Stdout)
	if err != nil {
		fatal(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: pdigit N")
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
