package pdigit

import (
	"bytes"
	"io"
	"testing"

	_ "embed"
)

func run(f func(r io.Reader, w io.Writer) error) []byte {
	r := bytes.NewReader(data)
	w := new(bytes.Buffer)
	err := f(r, w)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func runProcessor(n int) []byte {
	return run(Processor{GroupSpec: []int{n}, OutSep: SP}.Run)
}

func testF(t *testing.T, f func(int) []byte) {
	if res := f(3); !bytes.Equal(res, resultD3) {
		t.Errorf("mismatch\nhave:%s\nwant:%s\n", res, resultD3)
	}
	if res := f(4); !bytes.Equal(res, resultD4) {
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

var (
	//go:embed testdata/data1
	data []byte

	//go:embed testdata/data1r3
	resultD3 []byte

	//go:embed testdata/data1r4
	resultD4 []byte
)
