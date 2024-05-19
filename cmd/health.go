package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/idrac"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/angelhvargas/redfishcli/pkg/xclarity"
	"github.com/spf13/cobra"
)

var (
	drives     bool
	configPath string
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Return the Server RAID controllers health",
	Long: `This command returns the health status of RAID controllers for specified servers.
It can also return the health status of member drives if specified with the --drives flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		var servers []config.ServerConfig
		cfg, err := config.LoadConfigOrEnv(configPath, bmcType, bmcUsername, bmcPassword, bmcHost)
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

		// Marshal the health reports into JSON
		jsonData, err := json.Marshal(healthReports)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Print the JSON
		fmt.Println(string(jsonData))

		// Print errors if any
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
	// get the RAID controllers list (servers can have more than one)
	controllers, err := bmcClient.GetRAIDControllers()
	if err != nil {
		logger.Log.Error(err.Error())
		errorsCh <- err
		return
	}

	// get the data for each RAID controllers.
	for _, controller := range controllers {
		raidCtrldetails, err := bmcClient.GetRAIDControllerInfo(controller.ID)
		if err != nil {
			logger.Log.Error(err.Error())
			errorsCh <- err
			continue
		}

		healthReport := &model.RAIDHealthReport{
			ID:           raidCtrldetails.ID,
			Name:         raidCtrldetails.Name,
			HealthStatus: raidCtrldetails.Status.Health,
			State:        raidCtrldetails.Status.State,
			Drives:       []model.Drive{},
		}

		if drives {
			for _, driveRef := range raidCtrldetails.Drives {
				if len(driveRef.ID) > 0 {
					driveDetails, err := bmcClient.GetRAIDDriveDetails(driveRef.ID)
					if err != nil {
						logger.Log.Error(err.Error())
						errorsCh <- err
						continue
					}
					healthReport.Drives = append(healthReport.Drives, *driveDetails)
				}
			}
			healthReport.DrivesCount = int8(len(raidCtrldetails.Drives))
		}
		healthReportsCh <- healthReport
	}
}

func init() {
	raidCmd.AddCommand(healthCmd)
	healthCmd.PersistentFlags().BoolVarP(&drives, "drives", "", false, "return RAID controller member drives health")
}
