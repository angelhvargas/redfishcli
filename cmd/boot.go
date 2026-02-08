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
	"github.com/spf13/cobra"
)

// bootCmd represents the boot command
var bootCmd = &cobra.Command{
	Use:   "boot",
	Short: "Manage boot settings",
	Long:  `Manage server boot settings, including viewing boot info and setting the next boot device.`,
}

var bootStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get boot status and order",
	Run: func(cmd *cobra.Command, args []string) {
		runBootCommand(func(c client.ServerClient) error {
			info, err := c.GetBootInfo()
			if err != nil {
				return err
			}
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
			return nil
		})
	},
}

var bootSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set next boot device",
	Long:  `Set the next boot device. Common values: None, Pxe, Floppy, Cd, Usb, Hdd, BiosSetup, Utilities, Diags, UefiShell, UefiTarget`,
	Run: func(cmd *cobra.Command, args []string) {
		device, _ := cmd.Flags().GetString("device")
		if device == "" {
			logger.Log.Error("Device type is required")
			return
		}
		runBootCommand(func(c client.ServerClient) error {
			return c.SetBootOrder(device)
		})
	},
}

func runBootCommand(action func(client.ServerClient) error) {
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
	rootCmd.AddCommand(bootCmd)
	bootCmd.AddCommand(bootStatusCmd)
	bootCmd.AddCommand(bootSetCmd)
	bootSetCmd.Flags().StringP("device", "d", "", "Next boot device (e.g., Pxe, Hdd, Cd)")
}
