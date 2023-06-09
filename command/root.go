package command

import (
	"ant/utils/config"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   config.AppName,
	Short: fmt.Sprintf("%s is a project", config.AppName),
	Long:  fmt.Sprintf("%s is a project,this is long description", config.AppName),
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(qsCreateCmd)
	rootCmd.AddCommand(gameStartCmd)
	rootCmd.AddCommand(telegramCmd)
	// rootCmd.AddCommand(httpCmd)
}
