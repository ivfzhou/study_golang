//go:build linux

package bzip2_test

import (
	"bytes"
	"os"
	"testing"

	"gitee.com/ivfzhou/study_golang/invoke_c_code"
)

func TestNewBZip2Writer(t *testing.T) {
	file, err := os.ReadFile("../testdata/sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	buf := &bytes.Buffer{}
	bz := bzip2.NewBZip2Writer(buf)
	length, err := bz.Write(file)
	if err != nil {
		t.Fatal(err)
	}
	if length != len(file) {
		t.Fatal("length != len(file)")
	}
	err = bz.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("../testdata/sample.txt.bz2", buf.Bytes(), 0744)
	if err != nil {
		t.Fatal(err)
	}
}
