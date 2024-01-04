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
	"log"

	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/spf13/cobra"
)

// idracCmd represents the idrac command
var idracCmd = &cobra.Command{
	Use:   "idrac",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var servers []config.ServerConfig
		if cfgFile != "" {
			// Load configuration from file
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			servers = cfg.Servers
		} else {
			// Use the provided flags
			server := config.ServerConfig{
				Username: bmc_username,
				Password: bmc_password,
				Hostname: bmc_host,
			}
			servers = append(servers, server)
		}

		// Now you have a slice of ServerConfig (servers) to work with
		// Proceed with your application logic...
	},
}

var idracServerHealth = &cobra.Command{
	Use:   "health",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var servers []config.ServerConfig
		if cfgFile != "" {
			// Load configuration from file
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			servers = cfg.Servers
		} else {
			// Use the provided flags
			server := config.ServerConfig{
				Username: bmc_username,
				Password: bmc_password,
				Hostname: bmc_host,
			}
			servers = append(servers, server)
		}
		for _, server := range servers {
			fmt.Printf("Type: %s, Hostname: %s, Username: %s, Password: %s\n",
				server.Type, server.Hostname, server.Username, server.Password)
		}
		logger.Log.Info("Server Health called")

		// Now you have a slice of ServerConfig (servers) to work with
		// Proceed with your application logic...
	},
}

func init() {
	rootCmd.AddCommand(idracCmd)
	idracCmd.AddCommand(idracServerHealth)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// idracCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// idracCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
