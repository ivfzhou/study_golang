package get_windows_system_info

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func MemUnitParse(s string) (uint, error) {
	re := regexp.MustCompile("^(\\d+)([kKmMgG]?[bB]?)$")
	if !re.MatchString(s) {
		return 0, errors.New("unknown format")
	}
	arr := re.FindStringSubmatch(s)
	num, _ := strconv.ParseUint(arr[1], 10, 64)
	if len(arr) >= 3 {
		switch strings.ToLower(arr[2]) {
		case "b":
		case "k":
			num *= 1000
		case "m":
			num *= 1000 * 1000
		case "g":
			num *= 1000 * 1000 * 1000
		case "kb":
			num *= 1024
		case "mb":
			num *= 1024 * 1024
		case "gb":
			num *= 1024 * 1024 * 1024
		}
	}
	return uint(num), nil
}
