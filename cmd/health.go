/*
Copyright © 2024 Angel Vargas <angelvargas@outlook.es>

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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/angelhvargas/redfishcli/pkg/client"
	"github.com/angelhvargas/redfishcli/pkg/config"
	"github.com/angelhvargas/redfishcli/pkg/idrac"
	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/angelhvargas/redfishcli/pkg/model"
	"github.com/angelhvargas/redfishcli/pkg/xclarity"
	"github.com/spf13/cobra"
)

var (
	drives bool
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Return the Server RAID controllers health",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info(fmt.Sprintf("connecting to %s as user %s", bmc_host, bmc_username))
		var bmc_client client.ServerClient

		if bmc_type == "idrac" {
			logger.Log.Infoln("idrac client created...")
			bmc_client = idrac.NewClient(config.IDRACConfig{
				BMCConnConfig: config.BMCConnConfig{
					Hostname: bmc_host,
					Username: bmc_username,
					Password: bmc_password,
				},
			})

		} else {
			bmc_client = xclarity.NewClient(config.XClarityConfig{
				BMCConnConfig: config.BMCConnConfig{
					Hostname: bmc_host,
					Username: bmc_username,
					Password: bmc_password,
				},
			})
			logger.Log.Infoln("xclarity client created...")
		}

		controllers, err := bmc_client.GetRAIDControllers()
		if err != nil {
			logger.Log.Error(err.Error())
			panic(err.Error())
		}

		var healthReports []*model.RAIDHealthReport

		for _, controller := range controllers {
			details, err := bmc_client.GetRAIDControllerInfo(controller.ID)

			if err != nil {
				logger.Log.Error(err.Error())
				continue // Continue to the next controller if there's an error
			}

			healthReport := model.RAIDHealthReport{
				ID:           details.ID,
				Name:         details.Name,
				HealthStatus: details.Status.Health,
				State:        details.Status.State,
				Drives:       []model.Drive{},
			}

			if drives {

				for _, driveRef := range details.Drives {
					// Split the @odata.id string at the colon
					splitURL := strings.Split(driveRef.ID, ":")
					if len(splitURL) > 0 {
						// The first part of the split result is the valid URL
						validURL := splitURL[0]
						driveDetails, err := bmc_client.GetRAIDDriveDetails(validURL)
						if err != nil {
							logger.Log.Error(err.Error())
							continue // Continue to the next drive if there's an error
						}
						healthReport.Drives = append(healthReport.Drives, *driveDetails)
					}
				}
				healthReport.DrivesCount = int8(len(details.Drives))
			}

			healthReports = append(healthReports, &healthReport)
		}

		// Marshal the health reports into JSON
		jsonData, err := json.Marshal(healthReports)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Print the JSON
		fmt.Println(string(jsonData))

	},
}

func init() {
	raidCmd.AddCommand(healthCmd)
	healthCmd.PersistentFlags().BoolVarP(&drives, "drives", "", false, "return RAID controller member drives health")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
