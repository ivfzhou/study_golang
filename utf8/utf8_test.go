package utf8_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

/*
 * 1: 0xxx_xxxx 7 [0, 127] 0x0~0x7f
 * 2: 110x_xxxx 10xx_xxxx 11 [128, 2047] 0x80~0x7ff
 * 3: 1110_xxxx 10xx_xxxx 10xx_xxxx 16 [2048, 65535] 0x800~0xffff
 * 4: 1111_0xxx 10xx_xxxx 10xx_xxxx 10xx_xxxx 21 [65536, 2097151] 0x10000~0x1fffff
 *
 * 汉 {0b_1110_0110, 0b_1011_0001, 0b_1000_1001} 0x6c49
 */

func TestUTF8_1(t *testing.T) {
	res := strings.Builder{}
	for i, j := 0x0, 1; i <= 0x7f; i, j = i+1, j+1 {
		char := string([]byte{byte(i)})
		uc := strings.ReplaceAll(fmt.Sprintf("\\u%4X", i), " ", "0")
		hex := strings.ReplaceAll(fmt.Sprintf("\\x%2X", i), " ", "0")
		res.WriteString(uc + " " + hex + " " + char + "   ")
		if j%4 == 0 {
			res.WriteString("\n\n\n")
		}
	}
	err := os.WriteFile("./utf8_1.txt", []byte(res.String()), 0744)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTF8_2(t *testing.T) {
	res := strings.Builder{}
	for i, j := 0x80, 1; i <= 0x7ff; i, j = i+1, j+1 {
		b := strings.ReplaceAll(fmt.Sprintf("%11b", i), " ", "0")
		num, err := strconv.ParseUint("110"+b[:5], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		num2, err := strconv.ParseUint("10"+b[5:], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		char := string([]byte{byte(num), byte(num2)})
		uc := strings.ReplaceAll(fmt.Sprintf("\\u%4X", i), " ", "0")
		hex := strings.ReplaceAll(fmt.Sprintf("\\x%2X\\x%2X", num, num2), " ", "0")
		res.WriteString(uc + " " + hex + " " + char + "   ")
		if j%4 == 0 {
			res.WriteString("\n\n\n")
		}
	}
	err := os.WriteFile("./utf8_2.txt", []byte(res.String()), 0744)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTF8_3(t *testing.T) {
	res := strings.Builder{}
	for i, j := 0x800, 1; i <= 0xFFFF; i, j = i+1, j+1 {
		b := strings.ReplaceAll(fmt.Sprintf("%16b", i), " ", "0")
		num, err := strconv.ParseUint("1110"+b[:4], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		num1, err := strconv.ParseUint("10"+b[4:10], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		num2, err := strconv.ParseUint("10"+b[10:], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		char := string([]byte{byte(num), byte(num1), byte(num2)})
		uc := strings.ReplaceAll(fmt.Sprintf("\\u%4X", i), " ", "0")
		hex := fmt.Sprintf("\\x%2X\\x%2X\\x%2X", num, num1, num2)
		res.WriteString(uc + " " + hex + " " + char + "   ")
		if j%4 == 0 {
			res.WriteString("\n\n\n")
		}
	}
	err := os.WriteFile("./utf8_3.txt", []byte(res.String()), 0744)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTF8_4(t *testing.T) {
	res := strings.Builder{}
	for i, j := 0x10000, 1; i <= 0x1fffff; i, j = i+1, j+1 {
		b := strings.ReplaceAll(fmt.Sprintf("%21b", i), " ", "0")
		num, err := strconv.ParseUint("11110"+b[:3], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		num1, err := strconv.ParseUint("10"+b[3:9], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		num2, err := strconv.ParseUint("10"+b[9:15], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		num3, err := strconv.ParseUint("10"+b[15:], 2, 64)
		if err != nil {
			t.Fatal(err)
		}
		char := string([]byte{byte(num), byte(num1), byte(num2), byte(num3)})
		uc := strings.ReplaceAll(fmt.Sprintf("\\u%6X", i), " ", "0")
		hex := strings.ReplaceAll(fmt.Sprintf("\\x%2X\\x%2X\\x%2X\\x%2X", num, num1, num2, num3), " ", "0")
		res.WriteString(uc + " " + hex + " " + char + "   ")
		if j%4 == 0 {
			res.WriteString("\n\n\n")
		}
	}
	err := os.WriteFile("./utf8_4.txt", []byte(res.String()), 0744)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTF8All(t *testing.T) {
	res := strings.Builder{}
	res.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>`)
	for i, j := 0x0, 1; i <= 0x1fffff; i, j = i+1, j+1 {
		char := ""
		switch {
		case i <= 0x7f:
			char = string([]byte{byte(i)})
		case i <= 0x7ff && i >= 0x80:
			b := strings.ReplaceAll(fmt.Sprintf("%11b", i), " ", "0")
			num, err := strconv.ParseUint("110"+b[:5], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			num2, err := strconv.ParseUint("10"+b[5:], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			char = string([]byte{byte(num), byte(num2)})
		case i <= 0xffff && i >= 0x800:
			b := strings.ReplaceAll(fmt.Sprintf("%16b", i), " ", "0")
			num, err := strconv.ParseUint("1110"+b[:4], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			num1, err := strconv.ParseUint("10"+b[4:10], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			num2, err := strconv.ParseUint("10"+b[10:], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			char = string([]byte{byte(num), byte(num1), byte(num2)})
		case i <= 0x10000 && i >= 0x1fffff:
			b := strings.ReplaceAll(fmt.Sprintf("%21b", i), " ", "0")
			num, err := strconv.ParseUint("11110"+b[:3], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			num1, err := strconv.ParseUint("10"+b[3:9], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			num2, err := strconv.ParseUint("10"+b[9:15], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			num3, err := strconv.ParseUint("10"+b[15:], 2, 64)
			if err != nil {
				t.Fatal(err)
			}
			char = string([]byte{byte(num), byte(num1), byte(num2), byte(num3)})
		}
		res.WriteString(char + "  ")
	}
	res.WriteString(`</body>
</html>`)
	err := os.WriteFile("./utf8.html", []byte(res.String()), 0744)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUTF8AllValid(t *testing.T) {
	res := strings.Builder{}
	res.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>`)
	for i, j := uint32(0), 0; ; i++ {
		if utf8.ValidRune(rune(i)) {
			res.WriteString(fmt.Sprintf("%c ", i))
			j++
			if j%10 == 0 {
				res.WriteString("\n")
			}
		}
		if i == 1<<32-1 {
			break
		}
	}
	res.WriteString(`</body>
</html>`)
	err := os.WriteFile("./utf8_valid.html", []byte(res.String()), 0744)
	if err != nil {
		t.Fatal(err)
	}
}
