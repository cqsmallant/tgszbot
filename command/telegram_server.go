package command

import (
	"ant/plugin/telegram"

	"github.com/spf13/cobra"
)

var telegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "telegram服务",
	Long:  "telegram服务相关命令",
	Run: func(cmd *cobra.Command, args []string) {
		telegram.BotStart()
	},
}
