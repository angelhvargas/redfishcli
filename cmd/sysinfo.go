/*
Copyright © 2024 Angel Vargas <angelvargas@outlook.es>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// sysinfoCmd represents the sysinfo command
var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Get system information",
	Long:  `Retrieve detailed system information including BIOS version, serial number, model, and SKU.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigOrEnv(cfgFile, bmcType, bmcUsername, bmcPassword, bmcHost)
		if err != nil {
			logger.Log.Error(err.Error())
			os.Exit(1)
		}

		results := make([]*model.ServerInfo, 0)

		for _, server := range cfg.Servers {
			c, err := client.NewClient(server.Type, config.BMCConnConfig{
				Hostname: server.Hostname,
				Username: server.Username,
				Password: server.Password,
			})
			if err != nil {
				logger.Log.Errorf("Error creating client for server %s: %s", server.Hostname, err)
				continue
			}

			info, err := c.GetServerInfo()
			if err != nil {
				logger.Log.Errorf("Error getting sysinfo for server %s: %s", server.Hostname, err)
				continue
			}
			results = append(results, info)
		}

		printSysInfo(results)
	},
}

func printSysInfo(results []*model.ServerInfo) {
	switch output {
	case "json":
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	case "yaml":
		data, _ := yaml.Marshal(results)
		fmt.Println(string(data))
	default:
		// Simple text output
		for _, info := range results {
			fmt.Printf("ID: %s\n", info.ID)
			fmt.Printf("Manufacturer: %s\n", info.Manufacturer)
			fmt.Printf("Model: %s\n", info.Model)
			fmt.Printf("Serial Number: %s\n", info.SerialNumber)
			fmt.Printf("SKU: %s\n", info.SKU)
			fmt.Printf("BIOS Version: %s\n", info.BiosVersion)
			fmt.Printf("Power State: %s\n", info.PowerState)
			fmt.Printf("Health: %s\n", info.Status.Health)
			fmt.Println("--------------------------------------------------")
		}
	}
}

func init() {
	rootCmd.AddCommand(sysinfoCmd)
	sysinfoCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Output format (json, yaml, text)")
}
