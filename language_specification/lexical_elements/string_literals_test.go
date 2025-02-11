package lexical_elements_test

import (
	"fmt"
	"testing"
)

const (
	S1 = "\""
	S2 = "\u4f60\u597d\u4e16\u754c"
	S3 = "\U00004f60\U0000597d\U00004e16\U0000754c"
	S4 = "\xe4\xbd\xa0\xe5\xa5\xbd\xe4\xb8\x96\xe7\x95\x8c" // max=ff
	S5 = "\344\275\240\345\245\275\344\270\226\347\225\214" // max=377
	S6 = "\uffff"
	S7 = "\U0010FFFF"
)

func TestString(t *testing.T) {
	println(S1)
	println(S2)
	println(S3)
	println(S4)
	println(S5)
	println(S6)
	println(S7)
	fmt.Println([]byte("你好世界"))
}
