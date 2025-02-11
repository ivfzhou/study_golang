package properties_types_values_test

import "testing"

const (
	B  byte    = 'a'
	B1 byte    = 1.0
	I  int     = 0i
	F  float64 = 1 + 0i
)

func TestRepresentability(t *testing.T) {
	println(B)
	println(B1)
	println(I)
	println(F)
}
