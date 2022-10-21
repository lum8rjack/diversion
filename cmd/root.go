package cmd

import (
	"os"

	"github.com/lum8rjack/diversion/cmd/dump"
	"github.com/lum8rjack/diversion/cmd/patch"
	"github.com/lum8rjack/diversion/cmd/scan"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diversion",
	Short: "Suite of commands to use when performing Windows penetration testing",
	Long:  `Suite of commands to use when performing Windows penetration testing.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Adds the subcomments
func addSubcommandPalettes() {
	rootCmd.AddCommand(dump.DumpCmd)
	rootCmd.AddCommand(patch.PatchCmd)
	rootCmd.AddCommand(scan.ScanCmd)
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableFlagsInUseLine = true
	addSubcommandPalettes()
}
