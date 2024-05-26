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
	"gopkg.in/yaml.v3"
)

// controllersCmd represents the controllers command
var controllersCmd = &cobra.Command{
	Use:   "controllers",
	Short: "List and get details of storage controllers",
	Long: `This command lists and retrieves details of storage controllers for specified servers.

Usage:
  redfishcli storage controllers -t [controller-type] -u [username] -p [password] -n [hostname]

Options:
  -t, --bmc-type string   Controller type (idrac or xclarity) (default "idrac")
  -u, --username string   Username for the BMC
  -p, --password string   Password for the BMC
  -n, --host     string   Hostname or IP address of the server
  -o, --output   string   Output format (json, yaml, table) (default "json")

Example:
  List storage controllers of a Dell server with iDRAC:
    redfishcli storage controllers -t idrac -u root -p "your_password" -n 192.168.1.100 | jq

  List storage controllers of a Lenovo server with XClarity:
    redfishcli storage controllers -t xclarity -u admin -p "your_password" -n 192.168.1.101 | jq`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigOrEnv(cfgFile, bmcType, bmcUsername, bmcPassword, bmcHost)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		servers := cfg.Servers

		var wg sync.WaitGroup
		controllersReportsCh := make(chan *model.ControllersReport, len(servers))
		errorsCh := make(chan error, len(servers))

		for _, server := range servers {
			wg.Add(1)
			go processControllers(&wg, server, controllersReportsCh, errorsCh)
		}

		wg.Wait()
		close(controllersReportsCh)
		close(errorsCh)

		var controllersReports []*model.ControllersReport
		for report := range controllersReportsCh {
			controllersReports = append(controllersReports, report)
		}

		switch output {
		case "json":
			jsonData, err := json.Marshal(controllersReports)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			fmt.Println(string(jsonData))
		case "yaml":
			yamlData, err := yaml.Marshal(controllersReports)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
			fmt.Println(string(yamlData))
		case "table":
			printTable(controllersReports)
		default:
			logger.Log.Errorf("Unsupported output format: %s", output)
		}

		for err := range errorsCh {
			logger.Log.Error(err.Error())
		}
	},
}

func processControllers(wg *sync.WaitGroup, server config.ServerConfig, controllersReportsCh chan<- *model.ControllersReport, errorsCh chan<- error) {
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

	report, err := gatherControllersReport(bmcClient, server.Hostname)
	if err != nil {
		logger.Log.Error(err.Error())
		errorsCh <- err
		return
	}
	controllersReportsCh <- report
}

func gatherControllersReport(bmcClient client.ServerClient, hostname string) (*model.ControllersReport, error) {
	controllerConfig := &model.StorageControllerConfig{
		Type: "RAID",
	}
	controllers, err := bmcClient.GetStorageControllers(controllerConfig)
	if err != nil {
		return nil, err
	}

	report := &model.ControllersReport{
		Hostname:    hostname,
		Controllers: controllers,
	}

	return report, nil
}

func printTable(reports []*model.ControllersReport) {
	fmt.Printf("%-20s %-20s %-20s %-20s\n", "Hostname", "ID", "Name", "Status")
	for _, report := range reports {
		fmt.Printf("%-20s\n", report.Hostname)
		for _, controller := range report.Controllers {
			fmt.Printf("%-20s %-20s %-20s %-20s\n", "", controller.ID, controller.Name, controller.Status)
		}
	}
}

func init() {
	storageCmd.AddCommand(controllersCmd)
	controllersCmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "Output format (json, yaml, table)")
}
