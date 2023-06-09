package mq

import (
	"ant/model"
	"ant/plugin/telegram"
	"ant/utils/config"
	"ant/utils/log"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	QsStart   = "tgszbot:qs:start"
	QsLocking = "tgszbot:qs:locking"
	QsLocked  = "tgszbot:qs:locked"
	QsClosed  = "tgszbot:qs:closed"
)

func NewQsStartTask(qs *model.Qs) (*asynq.Task, error) {
	payload, err := json.Marshal(qs)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(QsStart, payload), nil
}

func QsStartHandle(ctx context.Context, t *asynq.Task) error {
	var qs model.Qs
	err := json.Unmarshal(t.Payload(), &qs)
	if err != nil {
		return err
	}
	qsLockingTime := qs.EndTime - time.Now().Unix() - 30
	defer func(qsLockingTime int64) {
		//前5秒已结束
		qsClosedQueue, _ := NewQsClosedTask(&qs)
		MClient.Enqueue(qsClosedQueue, asynq.ProcessIn(time.Duration(qsLockingTime+30)*time.Second))

		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}(qsLockingTime)

	if qs.ID > 0 && qs.Status == 0 {
		startGameStr := "<b>第<code>%s</code>期骰子游戏开始，请玩家开始投注，投注时间为%d秒。</b>"
		telegram.SendToBot(fmt.Sprintf(startGameStr, qs.Sn, config.QsStep))
		//更改状态
		qs.Status = 1
		model.EditQs(&qs)

		//前30封盘提醒
		qsLockingQueue, _ := NewQsLockingTask(&qs)
		MClient.Enqueue(qsLockingQueue, asynq.ProcessIn(time.Duration(qsLockingTime)*time.Second))

		//前10秒已封盘
		qsLockedQueue, _ := NewQsLockedTask(&qs)
		MClient.Enqueue(qsLockedQueue, asynq.ProcessIn(time.Duration(qsLockingTime+20)*time.Second))
	}

	return nil
}

func QsClosedHandle(ctx context.Context, t *asynq.Task) error {
	var qs model.Qs
	err := json.Unmarshal(t.Payload(), &qs)
	if err != nil {
		return err
	}
	defer func() {
		//新的一盘
		now := time.Now()
		newQs, _ := model.GetQsListByTime(now.Unix() + 10)
		qsStartQueue, _ := NewQsStartTask(newQs)
		MClient.Enqueue(qsStartQueue, asynq.ProcessIn(time.Duration(2)*time.Second))

		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()

	if qs.ID > 0 && qs.Status == 1 {
		//下注时间
		openFixGameStr := "<b>第<code>%s</code>期-开奖时间：%s\n\n—— —— ——封盘线—— —— —— \n\n</b>"
		orders := []string{"阳光  6杀 100  （赔率 1:4）", "阳光  6杀 100  （赔率 1:4）"}
		xzGameStrTemp := fmt.Sprintf(openFixGameStr, qs.Sn, "09:54:00")
		if len(orders) > 0 {
			xzGameStrTemp += "投注玩家\n"
			for _, item := range orders {
				xzGameStrTemp += item + "\n"
			}
			xzGameStrTemp += "\n——已封盘，线上下注全部有效——"
		} else {
			xzGameStrTemp += "无投注玩家"
		}
		telegram.SendToBot(xzGameStrTemp)
		time.Sleep(time.Second * 1)
		//筛子
		dice1Val := telegram.SendToDice()
		dice2Val := telegram.SendToDice()
		dice3Val := telegram.SendToDice()
		diceSum := dice1Val + dice2Val + dice3Val
		diceDx := "大"
		diceDs := "单"
		if diceSum < 11 {
			diceDx = "小"
		}
		if diceSum%2 == 0 {
			diceDs = "双"
		}

		//处理中奖结果
		//todo

		//开奖通知
		openGameStr := "<b>第<code>%s</code>期开奖结果：\n%s \n\n🎉🎉恭喜以下中奖玩家🎉🎉</b>"
		resGames := fmt.Sprintf("%d %d %d = %d %s %s", dice1Val, dice2Val, dice3Val, diceSum, diceDx, diceDs)
		telegram.SendToBot(fmt.Sprintf(openGameStr, qs.Sn, resGames))

		qs.Status = 2
		qs.Res = fmt.Sprintf("%d,%d,%d", dice1Val, dice2Val, dice3Val)
		model.EditQs(&qs)

		//----最近10期结果----20230318251期 5 6 3 = 14 大 双  20230318250期 1 4 1 = 6 小 双 对子
		// orderGameStr = "----最近10期结果----\n\n"
	}
	return nil
}

func NewQsClosedTask(qs *model.Qs) (*asynq.Task, error) {
	payload, err := json.Marshal(qs)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(QsClosed, payload), nil
}

// 封盘中
func NewQsLockingTask(qs *model.Qs) (*asynq.Task, error) {
	payload, err := json.Marshal(qs)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(QsLocking, payload), nil
}
func QsLockingHandle(ctx context.Context, t *asynq.Task) error {
	var qs model.Qs
	err := json.Unmarshal(t.Payload(), &qs)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()

	if qs.ID > 0 && qs.Status == 1 {
		//⏰封盘提醒 第20230318251期距离封盘还剩30秒,请玩家尽快投注！
		fpGameStr := "<b>⏰封盘提醒\n\n第<code>%s</code>期距离封盘还剩30秒,请玩家尽快投注！</b>"
		telegram.SendToBot(fmt.Sprintf(fpGameStr, qs.Sn))
	}
	return nil
}

// 已封盘
func NewQsLockedTask(qs *model.Qs) (*asynq.Task, error) {
	payload, err := json.Marshal(qs)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(QsLocked, payload), nil
}

func QsLockedHandle(ctx context.Context, t *asynq.Task) error {
	var qs model.Qs
	err := json.Unmarshal(t.Payload(), &qs)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()
	if qs.ID > 0 && qs.Status == 1 {
		//游戏已封盘，请勿投注，投注将视为无效投注！
		fpOkGameStr := "<b>第<code>%s</code>期已封盘，请勿投注，投注将视为无效投注！</b>"
		telegram.SendToBot(fmt.Sprintf(fpOkGameStr, qs.Sn))
	}
	return nil
}
