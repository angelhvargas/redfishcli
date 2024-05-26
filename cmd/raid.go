/*
Copyright © 2024 Angel Vargas <angelvargas@outlook.es>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// raidCmd represents the raid command
var raidCmd = &cobra.Command{
	Use:   "raid",
	Short: "Manage RAID health and details",
	Long: `The raid command provides functionalities to manage and query RAID health
and details for baremetal servers using Lenovo XClarity Controller or Dell iDRAC 7 or greater.

Usage:
  redfishcli storage raid [command]

Available Commands:
  health      Get the health status of RAID controllers
  details     Get details of RAID controllers

Examples:
  redfishcli storage raid health --drives -t xclarity -u admin -p "your_password" -n 192.168.1.101
  redfishcli storage raid details -t idrac -u root -p "your_password" -n 192.168.1.100`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("raid called")
	},
}

func init() {
	storageCmd.AddCommand(raidCmd)
}
