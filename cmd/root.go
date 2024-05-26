/*
Copyright © 2023 Angel Vargas <angelvargas@outlook.es>

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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	bmcUsername string
	bmcPassword string
	bmcHost     string
	bmcType     string
	output      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redfishcli",
	Short: "A tool to scan the health of baremetal servers using Redfish APIs",
	Long: `redfishcli is a command-line tool designed to scan the health of baremetal servers manufactured by Lenovo or Dell.
    
This tool provides a convenient way to monitor the health and status of your servers using Redfish APIs. It supports servers using Lenovo XClarity Controller or Dell iDRAC 7 or greater.

Features:
- Scan and report the health of baremetal servers.
- Support for Lenovo XClarity Controller.
- Support for Dell iDRAC 7 or greater.
- Retrieve RAID controller health status.
- Retrieve RAID drive details.
- Integration with Redfish APIs.

Installation:
- From Source:
  Ensure you have Go installed on your system, then run:
  go get github.com/angelhvargas/redfishcli
  cd $GOPATH/src/github.com/angelhvargas/redfishcli
  go install

- Using Precompiled Binaries:
  Precompiled binaries for various platforms are available on the releases page. Download the binary for your platform, extract it, and place it in a directory included in your system's PATH.

Usage:
- Basic Commands:
  Scan the RAID health of a server:
  redfishcli storage raid health --drives -t [controller-type] -u [username] -p [password] -n [hostname]

  -t: Controller type (idrac or xclarity).
  -u: Username for the BMC.
  -p: Password for the BMC.
  -n: Hostname or IP address of the server.

- Example Usage:
  Scan the RAID health of a Dell server with iDRAC:
  redfishcli storage raid health --drives -t idrac -u root -p "your_password" -n 192.168.1.100 | jq

  Scan the RAID health of a Lenovo server with XClarity:
  redfishcli storage raid health --drives -t xclarity -u admin -p "your_password" -n 192.168.1.101 | jq

- Configuration:
  You can create a configuration file to scan multiple servers without providing login parameters each time. By default, redfishcli looks for a configuration file at ~/.redfishcli/config.yaml.

  Example Configuration (config.yaml):
  servers:
    - type: "idrac"
      hostname: "192.168.1.100"
      username: "root"
      password: "your_password"
    - type: "xclarity"
      hostname: "192.168.1.101"
      username: "admin"
      password: "your_password"

  To use the configuration file, simply run:
  redfishcli storage raid health --drives
  redfishcli will automatically load the servers listed in the configuration file and scan their health.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command logic goes here
		if cfgFile == "" && (bmcUsername == "" || bmcPassword == "" || bmcHost == "") {
			cmd.Help() // Display help text
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.redfishcli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&bmcUsername, "username", "u", "", "username for server")
	rootCmd.PersistentFlags().StringVarP(&bmcPassword, "password", "p", "", "password for server")
	rootCmd.PersistentFlags().StringVarP(&bmcHost, "host", "n", "", "hostname of the server")
	rootCmd.PersistentFlags().StringVarP(&bmcType, "bmc-type", "t", "idrac", "BMC type (iDRAC or xClarity)")
	healthCmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "Output format (json, yaml, table)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".redfishcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".redfishcli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
