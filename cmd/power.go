/*
Copyright © 2024 Angel Vargas <angelvargas@outlook.es>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/spf13/cobra"
)

// powerCmd represents the power command
var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Manage server power state",
	Long: `Manage the power state of the server. 
Available commands are: status, on, off, restart.`,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current power state",
	Run: func(cmd *cobra.Command, args []string) {
		runPowerCommand(func(c client.ServerClient) error {
			state, err := c.GetPowerState()
			if err != nil {
				return err
			}
			fmt.Printf("Power State: %s\n", state)
			return nil
		})
	},
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Power on the server",
	Run: func(cmd *cobra.Command, args []string) {
		runPowerCommand(func(c client.ServerClient) error {
			return c.SetPowerState("On")
		})
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Power off the server (ForceOff by default)",
	Run: func(cmd *cobra.Command, args []string) {
		runPowerCommand(func(c client.ServerClient) error {
			return c.SetPowerState("ForceOff")
		})
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the server (GracefulRestart)",
	Run: func(cmd *cobra.Command, args []string) {
		runPowerCommand(func(c client.ServerClient) error {
			return c.Reboot()
		})
	},
}

func runPowerCommand(action func(client.ServerClient) error) {
	cfg, err := config.LoadConfigOrEnv(cfgFile, bmcType, bmcUsername, bmcPassword, bmcHost)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	for _, server := range cfg.Servers {
		fmt.Printf("Processing server: %s\n", server.Hostname)
		c, err := client.NewClient(server.Type, config.BMCConnConfig{
			Hostname: server.Hostname,
			Username: server.Username,
			Password: server.Password,
		})
		if err != nil {
			logger.Log.Errorf("Error creating client: %s", err)
			continue
		}

		if err := action(c); err != nil {
			logger.Log.Errorf("Error performing action: %s", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(powerCmd)
	powerCmd.AddCommand(statusCmd)
	powerCmd.AddCommand(onCmd)
	powerCmd.AddCommand(offCmd)
	powerCmd.AddCommand(restartCmd)
}
