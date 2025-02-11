package statements_test

import "testing"

type S struct{}

func (S) M(int) {}

func TestGo(t *testing.T) {
	go println("")

	go S.M(S{}, 1+1)
}
