package url_codec

import (
	"strconv"
	"strings"
)

const hexUpperCase = "0123456789ABCDEF"

/*
每个字节分成两部分表示（%FF）
*/

func Encode(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		if shouldEncode(b) {
			_, _ = buf.WriteString("%")
			_ = buf.WriteByte(hexUpperCase[b>>4])
			_ = buf.WriteByte(hexUpperCase[b&0b1111])
		} else {
			_ = buf.WriteByte(b)
		}
	}
	return buf.String()
}

func EscapeNonASCII(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b > 0x80 {
			_, _ = buf.WriteString("%")
			_, _ = buf.WriteString(strconv.FormatInt(int64(b), 16))
		} else {
			_ = buf.WriteByte(b)
		}
	}
	return buf.String()
}

func UnescapeNonASCII(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); {
		b := s[i]
		if b == '%' && i+2 < len(s) && canDecode(s[i+1]) && canDecode(s[i+2]) {
			encodedElem := string([]byte{s[i+1], s[i+2]})
			parseInt, _ := strconv.ParseUint(encodedElem, 16, 8)
			_ = buf.WriteByte(byte(parseInt))
			i += 3
		} else {
			_ = buf.WriteByte(b)
			i++
		}
	}
	return buf.String()
}

func Decode(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); {
		b := s[i]
		if b == '%' && i+2 < len(s) && canDecode(s[i+1]) && canDecode(s[i+2]) {
			_ = buf.WriteByte(decode(s[i+1])<<4 | decode(s[i+2]))
			i += 3
		} else {
			_ = buf.WriteByte(b)
			i++
		}
	}
	return buf.String()
}

func shouldEncode(char byte) bool {
	switch {
	case '0' <= char && char <= '9' ||
		'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z':
		return false
	default:
		return true
	}
}

func canDecode(char byte) bool {
	switch {
	case '0' <= char && char <= '9' ||
		'a' <= char && char <= 'f' ||
		'A' <= char && char <= 'F':
		return true
	default:
		return false
	}
}

func decode(char byte) byte {
	switch {
	case '0' <= char && char <= '9':
		return char - '0'
	case 'a' <= char && char <= 'f':
		return char - 'a' + 10
	case 'A' <= char && char <= 'F':
		return char - 'A' + 10
	default:
		return 0
	}
}
