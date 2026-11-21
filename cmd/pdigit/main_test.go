package main

import (
	"bytes"
	"io"
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
	// Note: all data is written before *fp is read, so it must fit in the OS pipe buffer
	// (64 kB on Linux).
	io.Copy(pw, b)
	pw.Close()
}
