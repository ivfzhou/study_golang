package strings

import "strings"

func Index(s1 string, s2 string) int {
	i, j := 0, 0
	next := indexNext(s2)
	for i < len(s1) && j < len(s2) {
		if j == 0 || s2[j] == s1[i] {
			i++
			j++
		} else {
			j = next[j]
		}
	}

	if j == len(s2) {
		return i - len(s2)
	}
	return 0
}

func Trim(s string, cutsets ...string) string {
	for {
		canBreak := true
		for _, cutset := range cutsets {
			if strings.HasPrefix(s, cutset) {
				canBreak = false
				s = s[len(cutset):]
			}
		}
		if canBreak {
			break
		}
	}
	for {
		canBreak := true
		for _, cutset := range cutsets {
			if strings.HasSuffix(s, cutset) {
				canBreak = false
				s = s[:len(s)-len(cutset)]
			}
		}
		if canBreak {
			break
		}
	}
	return s
}

func indexNext(s string) []int {
	res := make([]int, len(s))
	i, j := 1, 0
	for i < len(s)-1 {
		if j == 0 || s[i] == s[j] {
			i++
			j++
			if s[i] != s[j] {
				res[i] = j
			} else {
				res[i] = res[j]
			}
		} else {
			j = res[j]
		}
	}
	return res
}
