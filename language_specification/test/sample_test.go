package test_test

import (
	"testing"
	"unicode/utf8"
)

func TestT(t *testing.T) {

}

func BenchmarkB(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func ExampleFuncName() {
	// example code here.
}

// go test -v -fuzz ^FuzzF$ ./language_specification/test
func FuzzF(f *testing.F) {
	f.Add("ivfzhou")
	f.Fuzz(func(t *testing.T, name string) {
		rname := Reverse(name)
		rrname := Reverse(rname)
		if name != rrname {
			t.Error("name is not equaled", name, rrname)
		}
		if !utf8.ValidString(rname) {
			t.Error("name is not valid", rname)
		}
	})
}

func Reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
