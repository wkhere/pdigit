package pdigit

import (
	"bytes"
	"io"
	"testing"

	_ "embed"
)

type s = []int

func run(f func(r io.Reader, w io.Writer) error) []byte {
	r := bytes.NewReader(data)
	w := new(bytes.Buffer)
	err := f(r, w)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func runProcessor(spec []int) []byte {
	return run(Processor{GroupSpec: spec, OutSep: SP}.Run)
}

func TestProcessor(t *testing.T) {
	var tab = []struct {
		spec   []int
		result []byte
	}{
		{s{3}, resultD3},
		{s{4}, resultD4},
	}

	for i, tc := range tab {
		res := runProcessor(tc.spec)
		if !bytes.Equal(res, tc.result) {
			t.Errorf("tc#%d mismatch\nhave:%s\nwant:%s\n",
				i, res, tc.result)
		}
	}

}

func BenchmarkProcessor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runProcessor(s{3})
		runProcessor(s{4})
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
