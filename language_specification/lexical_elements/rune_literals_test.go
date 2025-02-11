package lexical_elements_test

import "testing"

const (
	R1 = '\a'
	R2 = '\b'
	R3 = '\f'
	R4 = '\n'
	R5 = '\r'
	R6 = '\t'
	R7 = '\v'
	R8 = '\\'
	R9 = '\''
)

func TestRune(t *testing.T) {
	println(string(R1))
	println(string(R2))
	println(string(R3))
	println(string(R4))
	println(string(R5))
	println(string(R6))
	println(string(R7))
	println(string(R8))
	println(string(R9))
}
