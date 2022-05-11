package cmd

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var configName, configValue string

var appConfigSetCmd = &cobra.Command{
	Use:     "set",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Set application config",
	Example: "civo app config set APP_NAME --name=foo --value=bar",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if configName == "PORT" {
			utility.Error("PORT is an immutable field picked by Civo.")
			os.Exit(1)
		}

		updatedConfig := civogo.EnvVar{
			Name:  configName,
			Value: configValue,
		}

		config := &civogo.UpdateApplicationRequest{
			Config: append(findApp.Config, updatedConfig),
		}

		app, err := client.UpdateApplication(findApp.ID, config)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": app.ID, "name": app.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Application %s's config has been updated.\n", utility.Green(app.Name))
		}
	},
}