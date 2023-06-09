package command

import (
	"ant/model"
	"ant/plugin/mq"
	"time"

	"github.com/spf13/cobra"
)

var gameStartCmd = &cobra.Command{
	Use:   "gameStart",
	Short: "游戏开始服务",
	Long:  "游戏开始生成",
	Run: func(cmd *cobra.Command, args []string) {
		gameStart()
	},
}

func gameStart() {
	now := time.Now()
	qs, _ := model.GetQsListByTime(now.Unix())

	qsStartQueue, _ := mq.NewQsStartTask(qs)
	mq.MClient.Enqueue(qsStartQueue)
}
