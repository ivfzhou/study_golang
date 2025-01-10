//go:build windows

package get_windows_system_info

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	mstrings "gitee.com/ivfzhou/study_golang/strings"
)

func GetCPU() (string, error) {
	op, err := exec.Command("cmd", "/Q", "/C", "TypePerf", "processor(_total)\\% processor time", "-sc", "1").Output()
	bs, _ := io.ReadAll(transform.NewReader(bytes.NewReader(op), simplifiedchinese.GBK.NewDecoder()))
	println(string(bs))
	if err != nil {
		return "", err
	}
	info := strings.Split(mstrings.Trim(string(bs), "\r", "\n", "\t", " "), "\r\n")
	for _, v := range info {
		arr := strings.Split(strings.TrimSpace(v), ",")
		if len(arr) != 2 {
			continue
		}
		percentage := strings.Trim(arr[1], `"`)
		_, err = strconv.ParseFloat(percentage, 64)
		if err != nil {
			continue
		}
		return percentage, nil
	}
	return "", nil
}

func GetMem() (total, free, vtotal, vfree uint, err error) {
	op, err := exec.Command("cmd", "/C", "systeminfo", "/FO", "CSV").Output()
	if err != nil {
		return
	}
	bs, _ := io.ReadAll(transform.NewReader(bytes.NewReader(op), simplifiedchinese.GBK.NewDecoder()))
	println(string(bs))
	if err != nil {
		return
	}

	arr := strings.Split(mstrings.Trim(string(bs), "\r", "\n", "\t", " "), "\r\n")
	if len(arr) < 2 {
		err = errors.New("unknown output")
		return
	}
	values := strings.Split(arr[1], `","`)
	if len(values) > 27 {
		totalS := strings.ReplaceAll(strings.ReplaceAll(values[22], ",", ""), " ", "")
		freeS := strings.ReplaceAll(strings.ReplaceAll(values[23], ",", ""), " ", "")
		vtotalS := strings.ReplaceAll(strings.ReplaceAll(values[24], ",", ""), " ", "")
		vfreeS := strings.ReplaceAll(strings.ReplaceAll(values[25], ",", ""), " ", "")

		total, _ = MemUnitParse(totalS)
		free, _ = MemUnitParse(freeS)
		vtotal, _ = MemUnitParse(vtotalS)
		vfree, _ = MemUnitParse(vfreeS)
	}

	return
}

func GetDisk() (total, free uint, err error) {
	op, err := exec.Command("cmd", "/C", "wmic", "volume").Output()
	if err != nil {
		return
	}
	bs, _ := io.ReadAll(transform.NewReader(bytes.NewReader(op), simplifiedchinese.GBK.NewDecoder()))
	println(string(bs))
	if err != nil {
		return
	}

	arr := strings.Split(mstrings.Trim(string(bs), "\r", "\n", "\t", " "), "\r\n")
	if len(arr) < 2 {
		err = errors.New("unknown output")
		return
	}

	titleIndex := make([]int, 1, 45)
	width := 0
	index := 0
	next := false
	for _, b := range arr[0] {
		if b != ' ' {
			if next {
				index++
				titleIndex = append(titleIndex, width)
				next = false
			}
		} else {
			next = true
		}
		width++
	}

	for _, v := range arr[1:] {
		if len(v) > titleIndex[len(titleIndex)-1] {
			letter := strings.ToUpper(strings.TrimSpace(v[titleIndex[14]:titleIndex[15]]))
			if letter == "C:" {
				capacity := strings.TrimSpace(v[titleIndex[5]:titleIndex[6]])
				freeSpace := strings.TrimSpace(v[titleIndex[20]:titleIndex[21]])
				total, _ = MemUnitParse(capacity)
				free, _ = MemUnitParse(freeSpace)
				return
			}
		}
	}

	return
}
