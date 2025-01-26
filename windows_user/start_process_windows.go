package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modwtsapi32                      = windows.NewLazySystemDLL("wtsapi32.dll")
	modkernel32                      = windows.NewLazySystemDLL("kernel32.dll")
	modadvapi32                      = windows.NewLazySystemDLL("advapi32.dll")
	moduserenv                       = windows.NewLazySystemDLL("userenv.dll")
	procWTSEnumerateSessionsW        = modwtsapi32.NewProc("WTSEnumerateSessionsW")
	procWTSGetActiveConsoleSessionId = modkernel32.NewProc("WTSGetActiveConsoleSessionId")
	procWTSQueryUserToken            = modwtsapi32.NewProc("WTSQueryUserToken")
	procDuplicateTokenEx             = modadvapi32.NewProc("DuplicateTokenEx")
	procCreateEnvironmentBlock       = moduserenv.NewProc("CreateEnvironmentBlock")
	procCreateProcessAsUser          = modadvapi32.NewProc("CreateProcessAsUserW")
	procGetTokenInformation          = modadvapi32.NewProc("GetTokenInformation")
	wtsQuerySessionInformationA      = modadvapi32.NewProc("WTSQuerySessionInformationA")
)

type WtsConnectStateClass int

type SecurityImpersonationLevel int

type TokenType int

type SW int

type WtsSessionInfo struct {
	SessionID      windows.Handle
	WinStationName *uint16
	State          WtsConnectStateClass
}

type TokenLinkedToken struct {
	LinkedToken windows.Token
}

const (
	WtsCurrentServerHandle uintptr = 0
)

const (
	WTSActive WtsConnectStateClass = iota
	WTSConnected
	WTSConnectQuery
	WTSShadow
	WTSDisconnected
	WTSIdle
	WTSListen
	WTSReset
	WTSDown
	WTSInit
)

const (
	SecurityAnonymous SecurityImpersonationLevel = iota
	SecurityIdentification
	SecurityImpersonation
	SecurityDelegation
)

const (
	TokenPrimary TokenType = iota + 1
	TokenImpersonazion
)

const (
	SWHide            SW = 0
	SWShowNormal         = 1
	SWNormal             = 1
	SWShowMinimized      = 2
	SWShowMaxImized      = 3
	SWMaxImize           = 3
	SWShowNoActivate     = 4
	SwShow               = 5
	SWMinimize           = 6
	SWShowMinnoactive    = 7
	SWShowna             = 8
	SWRestore            = 9
	SWShowDefault        = 10
	SWMax                = 1
)

const (
	CreateUnicodeEnvironment uint16 = 0x00000400
	CreateNoWindow                  = 0x08000000
	CreateNewConsole                = 0x00000010
)

// GetCurrentUserSessionID 获得当前系统活动的SessionID
func GetCurrentUserSessionID() (windows.Handle, error) {

	sessionList, err := WTSEnumerateSessions()
	if err != nil {
		return 0xFFFFFFFF, fmt.Errorf("get current user session token: %s", err)
	}
	for i := range sessionList {
		if sessionList[i].State == WTSActive {
			return sessionList[i].SessionID, nil
		}
	}
	if sessionId, _, err := procWTSGetActiveConsoleSessionId.Call(); sessionId == 0xFFFFFFFF {
		return 0xFFFFFFFF, fmt.Errorf("get current user session token: call native WTSGetActiveConsoleSessionId: %s", err)
	} else {
		return windows.Handle(sessionId), nil
	}
}

func GetUserSessionID(name string) (windows.Handle, error) {
	sessionList, err := WTSEnumerateSessions()
	if err != nil {
		return 0xFFFFFFFF, fmt.Errorf("get current user session token: %s", err)
	}
	for i := range sessionList {
		if sessionList[i].State == WTSActive {
			return sessionList[i].SessionID, nil
		}
	}
	if sessionId, _, err := procWTSGetActiveConsoleSessionId.Call(); sessionId == 0xFFFFFFFF {
		return 0xFFFFFFFF, fmt.Errorf("get current user session token: call native WTSGetActiveConsoleSessionId: %s", err)
	} else {
		return windows.Handle(sessionId), nil
	}
}

// WTSEnumerateSessions will call the native
// version for Windows and parse the result
// to a Golang friendly version
func WTSEnumerateSessions() ([]*WtsSessionInfo, error) {
	var (
		sessionInformation windows.Handle    = windows.Handle(0)
		sessionCount       int               = 0
		sessionList        []*WtsSessionInfo = make([]*WtsSessionInfo, 0)
	)
	if returnCode, _, err := procWTSEnumerateSessionsW.Call(WtsCurrentServerHandle, 0, 1, uintptr(unsafe.Pointer(&sessionInformation)), uintptr(unsafe.Pointer(&sessionCount))); returnCode == 0 {
		return nil, fmt.Errorf("call native WTSEnumerateSessionsW: %s", err)
	}
	structSize := unsafe.Sizeof(WtsSessionInfo{})
	current := uintptr(sessionInformation)
	for i := 0; i < sessionCount; i++ {
		sessionList = append(sessionList, (*WtsSessionInfo)(unsafe.Pointer(current)))
		current += structSize
	}
	return sessionList, nil
}

// DuplicateUserTokenFromSessionID will attempt
// to duplicate the user token for the user logged
// into the provided session ID
func DuplicateUserTokenFromSessionID(sessionId windows.Handle, runas bool) (windows.Token, error) {
	var (
		impersonationToken windows.Handle = 0
		userToken          windows.Token  = 0
	)

	if returnCode, _, err := procWTSQueryUserToken.Call(uintptr(sessionId), uintptr(unsafe.Pointer(&impersonationToken))); returnCode == 0 {
		return 0xFFFFFFFF, fmt.Errorf("call native WTSQueryUserToken: %s", err)
	}

	if returnCode, _, err := procDuplicateTokenEx.Call(uintptr(impersonationToken), 0, 0, uintptr(SecurityImpersonation), uintptr(TokenPrimary), uintptr(unsafe.Pointer(&userToken))); returnCode == 0 {
		return 0xFFFFFFFF, fmt.Errorf("call native DuplicateTokenEx: %s", err)
	}
	if runas {
		var admin TokenLinkedToken
		var dt uintptr = 0
		if returnCode, _, _ := procGetTokenInformation.Call(uintptr(impersonationToken), 19, uintptr(unsafe.Pointer(&admin)), uintptr(unsafe.Sizeof(admin)), uintptr(unsafe.Pointer(&dt))); returnCode != 0 {
			userToken = admin.LinkedToken
		}
	}
	if err := windows.CloseHandle(impersonationToken); err != nil {
		return 0xFFFFFFFF, fmt.Errorf("close windows handle used for token duplication: %s", err)
	}
	return userToken, nil
}

func StartProcessAsCurrentUser(appPath, cmdLine, workDir string, runas bool) error {
	var (
		sessionId windows.Handle
		userToken windows.Token
		envInfo   windows.Handle

		startupInfo windows.StartupInfo
		processInfo windows.ProcessInformation

		commandLine uintptr = 0
		workingDir  uintptr = 0

		err error
	)

	if sessionId, err = GetCurrentUserSessionID(); err != nil {
		return err
	}

	if userToken, err = DuplicateUserTokenFromSessionID(sessionId, runas); err != nil {
		return fmt.Errorf("get duplicate user token for current user session: %s", err)
	}

	if returnCode, _, err := procCreateEnvironmentBlock.Call(uintptr(unsafe.Pointer(&envInfo)), uintptr(userToken), 0); returnCode == 0 {
		return fmt.Errorf("create environment details for process: %s", err)
	}

	creationFlags := CreateUnicodeEnvironment | CreateNewConsole
	startupInfo.ShowWindow = SwShow
	startupInfo.Desktop = windows.StringToUTF16Ptr("winsta0\\default")

	if len(cmdLine) > 0 {
		commandLine = uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(cmdLine)))
	}
	if len(workDir) > 0 {
		workingDir = uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(workDir)))
	}
	if returnCode, _, err := procCreateProcessAsUser.Call(
		uintptr(userToken), uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(appPath))), commandLine, 0, 0, 0,
		uintptr(creationFlags), uintptr(envInfo), workingDir, uintptr(unsafe.Pointer(&startupInfo)), uintptr(unsafe.Pointer(&processInfo)),
	); returnCode == 0 {
		return fmt.Errorf("create process as user: %s", err)
	}
	return nil
}
