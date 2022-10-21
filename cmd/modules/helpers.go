package modules

import (
	"debug/pe"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Checks if the current process is 32bit
func IsWOW64Process() bool {
	var isWow64 bool

	handle := windows.CurrentProcess()
	err := windows.IsWow64Process(handle, &isWow64)
	if err != nil {
		log.Fatal("Error getting current process architecture")
	}

	return isWow64
}

func writeGoodBytes(b []byte, pn string, virtualoffset uint32, secname string, vsize uint32) {
	t, e := windows.LoadDLL(pn)
	if e != nil {
		fmt.Println("Error loading DLL from disk")
		os.Exit(1)
	}
	h := t.Handle
	dllBase := uintptr(h)

	dllOffset := uint(dllBase) + uint(virtualoffset)

	var old uint32
	e = windows.VirtualProtect(uintptr(dllOffset), uintptr(len(b)), windows.PAGE_READWRITE, &old)
	if e != nil {
		fmt.Printf("Error changing memory protections")
		os.Exit(1)
	}

	fmt.Println("Made memory map RW")

	for i := 0; i < len(b); i++ {
		loc := uintptr(dllOffset + uint(i))
		mem := (*[1]byte)(unsafe.Pointer(loc))
		(*mem)[0] = b[i]
	}

	fmt.Println("DLL overwritten")

	e = windows.VirtualProtect(uintptr(dllOffset), uintptr(len(b)), old, &old)
	if e != nil {
		fmt.Println("Error chaning permissions back to normal")
		os.Exit(1)
	}

	fmt.Println("Restored memory map permissions")
}

// Refreshes the provided DLL by reading the text section from disk and re-writing the dll in memeory
func ReloadDll(dll string) {
	fmt.Printf("Reloading DLL: %s\n", dll)

	df, e := ioutil.ReadFile(dll)
	if e != nil {
		fmt.Printf("Error reading DLL: %s\n", dll)
		os.Exit(1)
	}
	f, e := pe.Open(dll)
	if e != nil {
		fmt.Printf("Error opening DLL: %s\n", dll)
		os.Exit(1)
	}

	x := f.Section(".text")
	ddf := df[x.Offset:x.Size]
	writeGoodBytes(ddf, dll, x.VirtualAddress, x.Name, x.VirtualSize)
}

// Get PID of a provided process name
func GetPID(process string) uint32 {
	// Get PID
	h, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer windows.CloseHandle(h)

	// Find process pid
	PID := uint32(0)

	processEntrySize := unsafe.Sizeof(windows.ProcessEntry32{})
	p := windows.ProcessEntry32{Size: uint32(processEntrySize)}
	for {
		e := windows.Process32Next(h, &p)
		if e != nil {
			break
		}
		s := windows.UTF16ToString(p.ExeFile[:])
		if s == process {
			PID = p.ProcessID
		}
	}

	return PID
}

// Enable SePriv
func SePrivEnable(s string) error {
	var tokenHandle windows.Token
	thsHandle, err := windows.GetCurrentProcess()
	if err != nil {
		return err
	}
	windows.OpenProcessToken(
		//r, a, e := procOpenProcessToken.Call(
		thsHandle,                       //  HANDLE  ProcessHandle,
		windows.TOKEN_ADJUST_PRIVILEGES, //	DWORD   DesiredAccess,
		&tokenHandle,                    //	PHANDLE TokenHandle
	)
	var luid windows.LUID
	err = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr(s), &luid)
	if err != nil {
		log.Fatal("LookupPrivilegeValueW failed")
		return err
	}
	privs := windows.Tokenprivileges{}
	privs.PrivilegeCount = 1
	privs.Privileges[0].Luid = luid
	privs.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED
	err = windows.AdjustTokenPrivileges(tokenHandle, false, &privs, 0, nil, nil)
	if err != nil {
		log.Fatal("AdjustTokenPrivileges failed")
		return err
	}
	return nil
}
