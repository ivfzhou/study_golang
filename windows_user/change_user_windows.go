package main

import (
	"syscall"
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procWTSDisconnectSession = kernel32.NewProc("WTSDisconnectSession")
	procWTSConnectSession    = kernel32.NewProc("WTSConnectSession")
)

func Test1() {
	// Disconnect the current session
	_, _, err := procWTSDisconnectSession.Call(uintptr(syscall.Handle(0xFFFFFFFF)), uintptr(syscall.Handle(0)), 0)
	if err != nil {
		panic(err)
	}

	// Connect to session 1
	procWTSConnectSession.Call(uintptr(syscall.Handle(0xFFFFFFFF)), 1, 0, 0)
}
