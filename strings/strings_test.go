package strings_test

import (
	"testing"

	"gitee.com/ivfzhou/study_golang/strings"
)

func TestIndex(t *testing.T) {
	s1 := "acabaabaabcacaabc"
	s2 := "abaabc"
	println(strings.Index(s1, s2))

	s1 = "aaabaaaab"
	s2 = "aaaab"
	println(strings.Index(s1, s2))

}

func TestTrim(t *testing.T) {
	t.Log(strings.Trim(` 1234 11  1`, "1", " "))
}
