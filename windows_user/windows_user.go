package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	userInfos, err := GetUserInfo()
	if err != nil {
		log.Fatal(err)
	}
	for _, userInfo := range userInfos {
		fmt.Println(userInfo)
	}
}

// ===

type userInfo0 struct {
	Usri1_name *uint16
}

type localgroupUsersInfo0 struct {
	lgrui0_name *uint16
}

type userInfo2 struct {
	usri2_name           *uint16
	usri2_password       *uint16
	usri2_password_age   uint32
	usri2_priv           uint32
	usri2_home_dir       *uint16
	usri2_comment        *uint16
	usri2_flags          uint32
	usri2_script_path    *uint16
	usri2_auth_flags     uint32
	usri2_full_name      *uint16
	usri2_usr_comment    *uint16
	usri2_parms          *uint16
	usri2_workstations   *uint16
	usri2_last_logon     uint32
	usri2_last_logoff    uint32
	usri2_acct_expires   uint32
	usri2_max_storage    uint32
	usri2_units_per_week uint32
	usri2_logon_hours    *uint16
	usri2_bad_pw_count   uint32
	usri2_num_logons     uint32
	usri2_logon_server   *uint16
	usri2_country_code   uint32
	usri2_code_page      uint32
}

type UserInfo struct {
	Groups        []string
	UserType      string
	Name          string
	LastLoginTime uint32
}

func GetUserInfo() ([]*UserInfo, error) {
	netapi32 := syscall.NewLazyDLL("netapi32.dll")
	netUserEnum := netapi32.NewProc("NetUserEnum")
	netUserGetInfo := netapi32.NewProc("NetUserGetInfo")
	netUserGetLocalGroups := netapi32.NewProc("NetUserGetLocalGroups")
	netApiBufferFree := netapi32.NewProc("NetApiBufferFree")

	var (
		serverName                    [128]byte
		puserdata                     uintptr
		dwEntriesRead, dwTotalEntries uint32
	)
	bret, _, _ := netUserEnum.Call(
		uintptr(unsafe.Pointer(&serverName)),
		uintptr(0),
		uintptr(0x2),
		uintptr(unsafe.Pointer(&puserdata)),
		uintptr(128),
		uintptr(unsafe.Pointer(&dwEntriesRead)),
		uintptr(unsafe.Pointer(&dwTotalEntries)),
		uintptr(0),
	)
	if bret != 0 {
		return nil, fmt.Errorf("bret is %d", bret)
	}

	res := make([]*UserInfo, 0, dwEntriesRead)
	for i, iter := uint32(0), puserdata; i < dwEntriesRead; i++ {
		var (
			pgroupinfo                          uintptr
			groupEntriesread, groupTotalentries uint32
			userInfo                            UserInfo
		)
		data := (*userInfo0)(unsafe.Pointer(iter))
		bret, _, _ = netUserGetLocalGroups.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(data.Usri1_name)),
			uintptr(0), uintptr(0x1),
			uintptr(unsafe.Pointer(&pgroupinfo)),
			uintptr(0xFFFFFFFF),
			uintptr(unsafe.Pointer(&groupEntriesread)),
			uintptr(unsafe.Pointer(&groupTotalentries)),
		)
		if bret != 0 {
			iter = iter + unsafe.Sizeof(userInfo0{})
			continue
		}

		for j, ppgroupinfoIter := uint32(0), pgroupinfo; j < groupEntriesread; j++ {
			groupinfo := (*localgroupUsersInfo0)(unsafe.Pointer(ppgroupinfoIter))
			userInfo.Groups = append(userInfo.Groups, windows.UTF16PtrToString(groupinfo.lgrui0_name))
			ppgroupinfoIter = ppgroupinfoIter + unsafe.Sizeof(localgroupUsersInfo0{})
		}
		bret, _, _ = netApiBufferFree.Call(pgroupinfo)
		if bret != 0 {
			return nil, fmt.Errorf("bret is %d", bret)
		}

		var puserinfo uintptr
		bret, _, _ = netUserGetInfo.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(data.Usri1_name)),
			uintptr(2),
			uintptr(unsafe.Pointer(&puserinfo)),
		)

		if bret != 0 {
			iter = iter + unsafe.Sizeof(userInfo0{})
			continue
		}

		userdata := (*userInfo2)(unsafe.Pointer(puserinfo))

		userInfo.LastLoginTime = userdata.usri2_last_logon
		switch userdata.usri2_priv {
		case 0:
			userInfo.UserType = "GUEST"
		case 2:
			userInfo.UserType = "ADMIN"
		default:
			userInfo.UserType = "USER"
		}

		if userdata.usri2_name != nil {
			userInfo.Name = windows.UTF16PtrToString(userdata.usri2_name)
		}

		bret, _, _ = netApiBufferFree.Call(puserinfo)
		if bret != 0 {
			return nil, fmt.Errorf("bret is %d", bret)
		}
		iter = iter + unsafe.Sizeof(userInfo0{})
		res = append(res, &userInfo)
	}

	bret, _, _ = netApiBufferFree.Call(puserdata)
	if bret != 0 {
		return nil, fmt.Errorf("bret is %d", bret)
	}

	return res, nil
}
