package get_windows_system_info_test

import (
	"testing"

	wincmd "gitee.com/ivfzhou/study_golang/get_windows_system_info"
)

func TestMemUnitParse(t *testing.T) {
	t.Log(wincmd.MemUnitParse("1mb"))
	t.Log(wincmd.MemUnitParse("1gb"))
	t.Log(wincmd.MemUnitParse("1K"))
	t.Log(wincmd.MemUnitParse("1b"))
	t.Log(wincmd.MemUnitParse("mb"))
	t.Log(wincmd.MemUnitParse("1bb"))
	t.Log(wincmd.MemUnitParse("1bm"))
	t.Log(wincmd.MemUnitParse("1bm"))
}
