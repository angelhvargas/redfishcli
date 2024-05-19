package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/idrac"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/angelhvargas/redfishcli/pkg/xclarity"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	drives  bool
	output  string
	timeout time.Duration
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Return the Server RAID controllers health",
	Long: `This command returns the health status of RAID controllers for specified servers.
It can also return the health status of member drives if specified with the --drives flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		var servers []config.ServerConfig
		cfg, err := config.LoadConfigOrEnv(cfgFile, bmcType, bmcUsername, bmcPassword, bmcHost)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		servers = cfg.Servers

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
			printTable(healthReports)
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

	var bmcClient client.ServerClient
	switch server.Type {
	case "idrac":
		logger.Log.Infof("Creating iDRAC client for server %s", server.Hostname)
		bmcClient = idrac.NewClient(config.IDRACConfig{
			BMCConnConfig: config.BMCConnConfig{
				Hostname: server.Hostname,
				Username: server.Username,
				Password: server.Password,
			},
		})
	case "xclarity":
		logger.Log.Infof("Creating XClarity client for server %s", server.Hostname)
		bmcClient = xclarity.NewClient(config.XClarityConfig{
			BMCConnConfig: config.BMCConnConfig{
				Hostname: server.Hostname,
				Username: server.Username,
				Password: server.Password,
			},
		})
	default:
		err := fmt.Errorf("unsupported BMC type: %s", server.Type)
		logger.Log.Error(err.Error())
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
		return nil, fmt.Errorf("server is not powered on")
	}

	controllers, err := bmcClient.GetRAIDControllers()
	if err != nil {
		return nil, err
	}

	healthReport := &model.RAIDHealthReport{
		Hostname: hostname,
		Drives:   []model.Drive{},
	}

	for _, controller := range controllers {
		raidCtrldetails, err := bmcClient.GetRAIDControllerInfo(controller.ID)
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
					driveDetails, err := bmcClient.GetRAIDDriveDetails(driveRef.ID)
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

func printTable(reports []*model.RAIDHealthReport) {
	fmt.Printf("%-20s %-20s %-20s %-20s %-20s\n", "Hostname", "ID", "Name", "Health Status", "State")
	for _, report := range reports {
		fmt.Printf("%-20s %-20s %-20s %-20s %-20s\n", report.Hostname, report.ID, report.Name, report.HealthStatus, report.State)
		if drives {
			fmt.Println("Drives:")
			for _, drive := range report.Drives {
				fmt.Printf("  %-20s %-20s %-20s\n", drive.ID, drive.Status.Health, drive.Status.State)
			}
		}
	}
}

func init() {
	raidCmd.AddCommand(healthCmd)
	healthCmd.PersistentFlags().BoolVarP(&drives, "drives", "", false, "return RAID controller member drives health")
	healthCmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "Output format (json, yaml, table)")
	healthCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "", 60*time.Second, "Timeout duration for each server")
}
