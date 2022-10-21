package scan

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lum8rjack/Eagle/modules"
	"github.com/spf13/cobra"
)

var (
	inputList string
	ipAddress string
	open      bool
	output    string
	ports     string
	threads   int
	timeout   int
)

// ScanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Quick and portable TCP scanner",
	Long:  `Quick and portable TCP scanner. Provides some basic functionality like reading a list of hosts, saving the results to a file, and adjusting the timeout and number of hosts to scan concurrently`,
	Run: func(cmd *cobra.Command, args []string) {
		if inputList == "" && ipAddress == "" {
			cmd.Help()
			fmt.Println("\nYou must provide either a host to scan or list of hosts")
			os.Exit(1)
		}

		if ports == "" {
			cmd.Help()
			fmt.Println("\nYou must provide port(s) to scan")
			os.Exit(1)
		}

		// Setup the Logger
		logger := modules.Logger{
			SaveOutput: false,
			StartTime:  time.Now(),
			EndTime:    time.Now(),
		}

		if output != "" {
			logger.OutputFilename = output
			logger.SaveOutput = true
		}

		modules.SCANTIMEOUT = timeout

		portsString := strings.Split(ports, ",")
		portsInt, err := modules.CheckPorts(portsString)
		if err != nil {
			fmt.Println("Error parsing list of ports")
			os.Exit(1)
		}
		modules.NUMIPS = len(portsInt)

		var ipList []string

		// Get IPs from file or single host
		if inputList != "" {
			ipList = modules.ReadIPList(inputList)
		} else {
			ipList = modules.CheckCIDR(ipAddress)
		}

		// Setup number of go routines
		var wg sync.WaitGroup
		sem := make(chan int, threads)

		// Keep track of open ports
		openPorts := 0

		// Start scanning
		logger.Start()

		// Loop through IPs to scan
		for _, i := range ipList {
			wg.Add(1)
			sem <- 1

			go func(i string) {
				defer wg.Done()
				host := modules.InitialScan(i, portsInt)

				// Print the output from the scans
				for _, p := range host.Ports {
					var data string
					if open {
						if p.State == "Open" {
							data = fmt.Sprintf("%s:%d %s\n", host.Hostname, p.Port, p.State)
							openPorts++
						}
					} else {
						data = fmt.Sprintf("%s:%d %s\n", host.Hostname, p.Port, p.State)
					}
					logger.WriteToFile(data)
				}
				<-sem
			}(i)
		}

		// Wait for scans to be done
		wg.Wait()
		close(sem)

		// Scanning completed
		logger.Stop()

		// Log how many ports were open
		fmt.Printf("Number of open ports: %v\n", openPorts)
	},
}

func init() {
	ScanCmd.DisableFlagsInUseLine = true
	ScanCmd.Flags().StringVarP(&inputList, "file", "f", "", "Input file containing a list of hosts to scan")
	ScanCmd.Flags().StringVarP(&ipAddress, "ip", "i", "", "Single host or CIDR address to scan")
	ScanCmd.Flags().BoolVarP(&open, "open", "o", true, "Only output open ports")
	ScanCmd.Flags().StringVarP(&output, "save", "s", "", "File to save the results")
	ScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (comma separated)")
	ScanCmd.Flags().IntVarP(&threads, "threads", "d", 20, "Number of hosts to scan concurrently")
	ScanCmd.Flags().IntVarP(&timeout, "timeout", "t", 3, "Timeout in seconds for each scanned port")
}
