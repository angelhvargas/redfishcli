package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/angelhvargas/redfishcli/pkg/tableprinter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	drives  bool
	timeout time.Duration
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Return the Server RAID controllers health",
	Long: `This command returns the health status of RAID controllers for specified servers.

It can also return the health status of member drives if specified with the --drives flag.

Usage:
  redfishcli storage raid health --drives -t [controller-type] -u [username] -p [password] -n [hostname]

Options:
  --drives       Include health status of RAID member drives
  -t, --bmc-type string   Controller type (idrac or xclarity) (default "idrac")
  -u, --username string   Username for the BMC
  -p, --password string   Password for the BMC
  -n, --host     string   Hostname or IP address of the server
  -o, --output   string   Output format (json, yaml, table) (default "json")

Example:
  Scan the RAID health of a Dell server with iDRAC:
    redfishcli storage raid health --drives -t idrac -u root -p "your_password" -n 192.168.1.100 | jq

  Scan the RAID health of a Lenovo server with XClarity:
    redfishcli storage raid health --drives -t xclarity -u admin -p "your_password" -n 192.168.1.101 | jq

Configuration:
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
		cfg, err := config.LoadConfigOrEnv(cfgFile, bmcType, bmcUsername, bmcPassword, bmcHost)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		servers := cfg.Servers

		var wg sync.WaitGroup
		healthReportsCh := make(chan *model.RAIDHealthReport, len(servers))
		errorsCh := make(chan error, len(servers))

		for _, server := range servers {
			wg.Add(1)
			go processServer(&wg, server, healthReportsCh, errorsCh)
		}

		wg.Wait()
		close(healthReportsCh)
		close(errorsCh)

		var healthReports []*model.RAIDHealthReport
		for report := range healthReportsCh {
			healthReports = append(healthReports, report)
		}

		switch output {
		case "json":
			jsonData, err := json.Marshal(healthReports)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			fmt.Println(string(jsonData))
		case "yaml":
			yamlData, err := yaml.Marshal(healthReports)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			fmt.Println(string(yamlData))
		case "table":
			headers := []string{"Hostname", "ID", "Name", "Health Status", "State"}
			fields := []string{"Hostname", "ID", "Name", "HealthStatus", "State"}
			nestedConfig := map[string]tableprinter.NestedTableConfig{
				"Drives": {
					Headers: []string{"Drive ID", "Health", "State"},
					Fields:  []string{"ID", "Status.Health", "Status.State"},
				},
			}
			tableprinter.PrintTable(healthReports, headers, fields, nestedConfig, 0)
		default:
			logger.Log.Errorf("Unsupported output format: %s", output)
		}

		for err := range errorsCh {
			logger.Log.Error(err.Error())
		}
	},
}

func processServer(wg *sync.WaitGroup, server config.ServerConfig, healthReportsCh chan<- *model.RAIDHealthReport, errorsCh chan<- error) {
	defer wg.Done()

	// Create client using the registry
	bmcClient, err := client.NewClient(server.Type, config.BMCConnConfig{
		Hostname: server.Hostname,
		Username: server.Username,
		Password: server.Password,
	})
	if err != nil {
		logger.Log.Errorf("Error creating client for server %s: %s", server.Hostname, err)
		errorsCh <- err
		return
	}

	report, err := gatherHealthReport(bmcClient, server.Hostname)
	if err != nil {
		logger.Log.Error(err.Error())
		errorsCh <- err
		// If there is an error, create a report with "unknown" state
		report = &model.RAIDHealthReport{
			Hostname:     server.Hostname,
			State:        "unknown",
			HealthStatus: "unknown",
		}
	}
	healthReportsCh <- report
}

func gatherHealthReport(bmcClient client.ServerClient, hostname string) (*model.RAIDHealthReport, error) {
	serverStatus, err := bmcClient.GetServerInfo()
	if err != nil {
		return nil, err
	}

	if serverStatus.PowerState != "On" {
		return nil, fmt.Errorf("host %s: server is not powered on", hostname)
	}
	config := &model.StorageControllerConfig{
		Type: "RAID",
	}
	controllers, err := bmcClient.GetStorageControllers(config)
	if err != nil {
		return nil, err
	}

	healthReport := &model.RAIDHealthReport{
		Hostname: hostname,
		Drives:   []model.Drive{},
	}

	for _, controller := range controllers {
		raidCtrldetails, err := bmcClient.GetStorageControllerInfo(controller.ID)
		if err != nil {
			return nil, err
		}

		healthReport.ID = raidCtrldetails.ID
		healthReport.Name = raidCtrldetails.Name
		healthReport.HealthStatus = raidCtrldetails.Status.Health
		healthReport.State = raidCtrldetails.Status.State

		if drives {
			for _, driveRef := range raidCtrldetails.Drives {
				if len(driveRef.ID) > 0 {
					driveDetails, err := bmcClient.GetStorageDriveDetails(driveRef.ID)
					if err != nil {
						return nil, err
					}
					healthReport.Drives = append(healthReport.Drives, *driveDetails)
				}
			}
			healthReport.DrivesCount = int8(len(raidCtrldetails.Drives))
		}
	}

	return healthReport, nil
}

func init() {
	raidCmd.AddCommand(healthCmd)
	healthCmd.PersistentFlags().BoolVarP(&drives, "drives", "", false, "return RAID controller member drives health")
	healthCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "", 60*time.Second, "Timeout duration for each server")
}
