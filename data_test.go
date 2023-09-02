package pdigit

import (
	"bytes"
	"testing"

	_ "embed"
)

func runData(data []byte, spec []int) []byte {
	r := bytes.NewReader(data)
	b := new(bytes.Buffer)
	p := Proc{GroupSpec: spec, OutSep: SP}
	err := p.Run(r, b)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func TestData1(t *testing.T) {
	var tab = []struct {
		spec   []int
		result []byte
	}{
		{s{3}, data1r3},
		{s{4}, data1r4},
		{s{2, 4}, data1r24},
	}

	for i, tc := range tab {
		res := runData(data1, tc.spec)
		if !bytes.Equal(res, tc.result) {
			t.Errorf("tc#%d mismatch\nhave:%s\nwant:%s\n",
				i, res, tc.result)
		}
	}

}

func BenchmarkData1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runData(data1, s{3})
		runData(data1, s{4})
		runData(data1, s{1, 2})
		runData(data1, s{2, 4})
	}
}

var (
	//go:embed testdata/data1
	data1 []byte

	//go:embed testdata/data1r3
	data1r3 []byte

	//go:embed testdata/data1r4
	data1r4 []byte

	//go:embed testdata/data1r2,4
	data1r24 []byte
)
