package pdigit

import (
	"bytes"
	"io"
	"testing"
)

func run(f func(r io.Reader, w io.Writer) error) string {
	r := bytes.NewBufferString(data)
	w := bytes.NewBuffer(nil)
	err := f(r, w)
	if err != nil {
		panic(err)
	}
	return w.String()
}

func runProcessor(n int) string {
	return run(Processor{NDigits: n, OutSep: SP}.Run)
}

func testF(t *testing.T, f func(int) string) {
	if res := f(3); res != resultD3 {
		t.Errorf("mismatch\nhave:%s\nwant:%s\n", res, resultD3)
	}
	if res := f(4); res != resultD4 {
		t.Errorf("mismatch\nhave:%s\nwant:%s\n", res, resultD4)
	}
}

func TestProcessor(t *testing.T) {
	testF(t, runProcessor)
}

func BenchmarkProcessor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runProcessor(3)
		runProcessor(4)
	}
}

const data = `

abc
a b
a 12345 b
a12345 b
a 12345b
a12345b
1
12
123
1234
12345
123456
1234567
12345678
123456789
1234567890
1 12345
1 a 12345
1234 12345
1234 x 12345
12345 123456
12345 x 123456
`

const resultD3 = `

abc
a b
a 12 345 b
a12 345 b
a 12345b
a12345b
1
12
123
1 234
12 345
123 456
1 234 567
12 345 678
123 456 789
1 234 567 890
1 12 345
1 a 12 345
1 234 12 345
1 234 x 12 345
12 345 123 456
12 345 x 123 456
`

const resultD4 = `

abc
a b
a 1 2345 b
a1 2345 b
a 12345b
a12345b
1
12
123
1234
1 2345
12 3456
123 4567
1234 5678
1 2345 6789
12 3456 7890
1 1 2345
1 a 1 2345
1234 1 2345
1234 x 1 2345
1 2345 12 3456
1 2345 x 12 3456
`
