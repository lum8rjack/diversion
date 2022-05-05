package modules

import (
	"log"

	"golang.org/x/sys/windows"
)

func IsWOW64Process() bool {
	var isWow64 bool

	handle := windows.CurrentProcess()
	err := windows.IsWow64Process(handle, &isWow64)
	if err != nil {
		log.Fatal("Error getting current process architecture")
	}

	return isWow64
}
