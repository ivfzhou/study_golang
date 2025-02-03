package expressions_test

import (
	"fmt"
	"testing"
)

func TestLabel(t *testing.T) {
	l1 := "1"
l1:
	{
		fmt.Println(l1)
		goto l2
	}

l3:
	for {
		break
	}

l2:
	for {
		func() {
			// goto l1
		}()
		//continue l3
		goto l3
	}
	goto l1
	goto l3
}
