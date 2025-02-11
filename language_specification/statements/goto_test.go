package statements_test

import "testing"

func TestGoto(t *testing.T) {
	//{
l1:
	{
		println("ok")
	}

	//}

	goto l1
}
