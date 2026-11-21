package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

var bufWrite = bytes.NewBufferString

func Example() {
	os.Args = append(os.Args[:1], "2")

	b := bufWrite("00112233")

	feed(&os.Stdin, b)
	main()

	// Output:
	// 00 11 22 33
}

func Example_withSeparator() {
	os.Args = append(os.Args[:1], "-s|", "2")

	b := bufWrite("00112233")

	feed(&os.Stdin, b)
	main()

	// Output:
	// 00|11|22|33
}

func Example_withSpacedSeparator() {
	os.Args = append(os.Args[:1], "-s", "|", "2")

	b := bufWrite("00112233")

	feed(&os.Stdin, b)
	main()

	// Output:
	// 00|11|22|33
}

func Example_withDashDash() {
	os.Args = append(os.Args[:1], "--", "2")

	b := bufWrite("00112233")

	feed(&os.Stdin, b)
	main()

	// Output:
	// 00 11 22 33
}

func feed(fp **os.File, b io.Reader) {
	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	*fp = pr
	ioutil.ReadAll(io.TeeReader(b, pw))
	pw.Close()
}
