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

// eventlogCmd represents the eventlog command
var eventlogCmd = &cobra.Command{
	Use:   "eventlog",
	Short: "Get system event logs",
	Long:  `Retrieve the System Event Log (SEL).`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigOrEnv(cfgFile, bmcType, bmcUsername, bmcPassword, bmcHost)
		if err != nil {
			logger.Log.Error(err.Error())
			os.Exit(1)
		}

		for _, server := range cfg.Servers {
			fmt.Printf("--- Event Logs for %s ---\n", server.Hostname)
			c, err := client.NewClient(server.Type, config.BMCConnConfig{
				Hostname: server.Hostname,
				Username: server.Username,
				Password: server.Password,
			})
			if err != nil {
				logger.Log.Errorf("Error creating client: %s", err)
				continue
			}

			logs, err := c.GetSystemEventLog()
			if err != nil {
				logger.Log.Errorf("Error getting event logs: %s", err)
				continue
			}

			printEventLogs(logs)
		}
	},
}

func printEventLogs(logs []model.EventLogEntry) {
	switch output {
	case "json":
		data, _ := json.MarshalIndent(logs, "", "  ")
		fmt.Println(string(data))
	case "yaml":
		data, _ := yaml.Marshal(logs)
		fmt.Println(string(data))
	default:
		// Simple text output
		for _, entry := range logs {
			fmt.Printf("[%s] %s: %s (Severity: %s)\n", entry.Created, entry.EntryType, entry.Message, entry.Severity)
		}
	}
}

func init() {
	rootCmd.AddCommand(eventlogCmd)
	eventlogCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Output format (json, yaml, text)")
}
