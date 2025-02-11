package statements_test

import (
	"fmt"
	"testing"
)

func TestReturn(t *testing.T) {
	fmt.Println(NoErr())
}

func NoErr() (s string, _ error) {
	s = "ok"

	/* if s := "ok1"; s != "" {
		return // 这里是 return 外面的 s 还是里面的s，未知，无法编译。
	} */
	return
}
