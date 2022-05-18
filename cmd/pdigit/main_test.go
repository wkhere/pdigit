package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func ExampleMain() {
	os.Args = os.Args[:1]
	os.Args = append(os.Args, "2")

	b := new(bytes.Buffer)
	b.WriteString("00112233")

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
