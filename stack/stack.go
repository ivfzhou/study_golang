package main

import (
	"fmt"
	"runtime"
	"strings"
)

func main() {
	SampleUse()
}

func SampleUse() {
	defer func() {
		p := recover()
		if p != nil {
			callers := GetStackCallers("certs-master")
			fmt.Println(p, callers)
		}
	}()

	FnA()
}

func FnA() {
	FnB()
}

func FnB() {
	FnC()
}

func FnC() {
	FnD()
}

func FnD() {
	FnE()
}

func FnE() {
	panic("i am panic")
}

func GetStackCallers(flag string) string {
	callers := make([]uintptr, 32)
	n := runtime.Callers(3, callers)
	callers = callers[:n]
	frames := runtime.CallersFrames(callers)
	elems := make([]string, 0, len(callers))
	for {
		frame, more := frames.Next()
		index := strings.LastIndex(frame.File, flag)
		if index != -1 {
			elems = append(elems, fmt.Sprintf("%s:%d", frame.File[index+len(flag):], frame.Line))
		}

		if !more {
			break
		}
	}
	return strings.Join(elems, ", ")
}
