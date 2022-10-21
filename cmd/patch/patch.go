package patch

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	pid    int
	ppid   bool
	method string
)

// PatchCmd represents the patch command
var PatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Patch the specified function",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Get PPID
		if ppid {
			fmt.Printf("PPID: %d\n", os.Getppid())
			os.Exit(0)
		}

		// Make sure user supplied PID
		if pid == 0 {
			cmd.Help()
			fmt.Println("\n[!] Required flag \"pid\" not set")
			os.Exit(0)
		}

		// Parse method and start patching
		m := strings.ToLower(method)
		if m == "amsi" {
			err := PatchAmsi(pid)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Successfully patched AMSI.AmsiOpenSession in remote process with PID: %d\n", pid)
		} else if m == "etw" {
			err := PatchETW(pid)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Successfully patched NTDLL.EtwEventWrite in remote process with PID: %d\n", pid)
		} else {
			fmt.Printf("Invalid method provided: %s\n", method)
		}
	},
}

func init() {
	PatchCmd.DisableFlagsInUseLine = true
	PatchCmd.Flags().StringVarP(&method, "method", "m", "amsi", "Evasion method: amsi,etw")
	PatchCmd.Flags().IntVarP(&pid, "pid", "p", 0, "PID of the process to inject into")
	PatchCmd.Flags().BoolVarP(&ppid, "ppid", "i", false, "Get the PPID of current process")
}
