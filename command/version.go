package command

import (
	"ant/utils/config"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s", config.AppName),
	Long:  fmt.Sprintf("All software has versions. This is %s", config.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		print(config.GetAppVersion())
	},
}
