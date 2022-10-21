package dump

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	pid     int
	outfile string
)

// DumpCmd represents the dump command
var DumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump process memory using MiniDumpWriteDump",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Make sure user supplied PID
		if pid == 0 {
			cmd.Help()
			fmt.Println("\n[!] Required flag \"pid\" not set")
			os.Exit(0)
		}
		// Make sure user supplied outfile
		if outfile == "" {
			cmd.Help()
			fmt.Println("\n[!] Required flag \"outfile\" not set")
			os.Exit(0)
		}

		// Dump memory
		minidump(pid, outfile)
	},
}

func init() {
	DumpCmd.DisableFlagsInUseLine = true
	DumpCmd.Flags().IntVarP(&pid, "pid", "p", 0, "PID of the process to dump")
	DumpCmd.Flags().StringVarP(&outfile, "outfile", "o", "", "File to save the memory dump to")
}
