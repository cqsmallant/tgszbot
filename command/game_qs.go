package command

import (
	"ant/model"
	"ant/utils/config"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var qsCreateCmd = &cobra.Command{
	Use:   "qscreate",
	Short: "期数服务",
	Long:  "期数服务生成",
	Run: func(cmd *cobra.Command, args []string) {
		qsCreated()
	},
}

func qsCreated() {
	now := time.Now()
	qsStep := config.QsStep
	totalQs := 24 * 3600 / qsStep

	//判断是否存在，
	qs, err := model.GetQsListByTime(now.Unix())
	if err != nil {
		panic(err)
	}

	qsBeanList := []model.Qs{}
	curDay := now.Format("20060102")
	toTime, _ := time.ParseInLocation("20060102", curDay, time.Local)
	//判断今天是否存在
	if qs.ID < 1 {
		//生成今天的数据
		for i := 1; i <= totalQs; i++ {
			tempQs := curDay
			if i < 100 {
				tempQs += "0"
			}
			if i < 10 {
				tempQs += "0"
			}
			tempQs = fmt.Sprintf(tempQs+"%d", i)
			qsBean := model.Qs{
				Sn:        tempQs,
				BeginTime: toTime.Add(time.Second * time.Duration((i-1)*qsStep)).Unix(),
				EndTime:   toTime.Add(time.Second * time.Duration(i*qsStep)).Unix(),
				Status:    0,
			}
			qsBean.CreateTime = now.Unix()
			qsBeanList = append(qsBeanList, qsBean)
		}
		model.AddQsInBatches(&qsBeanList)
	}

	//判断明天是否存在
	tomorrowCurTime := now.Add(time.Hour * 24)
	tomorrowDay := tomorrowCurTime.Format("20060102")
	tomorrowTime2, _ := time.ParseInLocation("20060102", tomorrowDay, time.Local)
	qs, err = model.GetQsListByTime(tomorrowCurTime.Unix())
	if err != nil {
		panic(err)
	}
	if qs.ID < 1 {
		qsBeanList = []model.Qs{}
		//生成明天的数据
		for i := 1; i <= totalQs; i++ {
			tempQs := tomorrowDay
			if i < 100 {
				tempQs += "0"
			}
			if i < 10 {
				tempQs += "0"
			}
			tempQs = fmt.Sprintf(tempQs+"%d", i)

			qsBean := model.Qs{
				Sn:        tempQs,
				BeginTime: tomorrowTime2.Add(time.Second * time.Duration((i-1)*qsStep)).Unix(),
				EndTime:   tomorrowTime2.Add(time.Second * time.Duration(i*qsStep)).Unix(),
				Status:    0,
			}
			qsBean.CreateTime = now.Unix()
			qsBeanList = append(qsBeanList, qsBean)
		}
		model.AddQsInBatches(&qsBeanList)
	}
}
