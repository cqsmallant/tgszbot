package mq

import (
	"ant/model"
	"ant/plugin/telegram"
	"ant/utils/config"
	"ant/utils/constant"
	"ant/utils/dao"
	"ant/utils/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

const (
	QsStart   = "tgszbot:qs:start"
	QsLocking = "tgszbot:qs:locking"
	QsLocked  = "tgszbot:qs:locked"
	QsClosed  = "tgszbot:qs:closed"
)

var mutexLock sync.Mutex

func NewQsStartTask(qs *model.Qs) (*asynq.Task, error) {
	payload, err := json.Marshal(qs)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(QsStart, payload), nil
}

func QsStartHandle(ctx context.Context, t *asynq.Task) error {
	mutexLock.Lock()
	var qs model.Qs
	err := json.Unmarshal(t.Payload(), &qs)
	if err != nil {
		return err
	}
	defer func() {
		mutexLock.Unlock()
		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()
	if qs.ID > 0 && qs.Status == 0 {
		startGameStr := "<b>第<code>%s</code>期骰子游戏开始，请玩家开始投注，投注时间为%d秒。</b>"
		telegram.SendToBot(fmt.Sprintf(startGameStr, qs.Sn, config.QsStep))
		//更改状态
		qs.Status = 1
		model.EditQs(&qs)

		qsTime := qs.EndTime - time.Now().Unix()
		//设置缓存
		payload, _ := json.Marshal(qs)
		dao.Rdb.Set(ctx, constant.CacheQsNow, payload, time.Duration(qsTime)*time.Second)
		//前30封盘提醒
		qsLockingQueue, _ := NewQsLockingTask(&qs)
		MClient.Enqueue(qsLockingQueue, asynq.ProcessIn(time.Duration(qsTime-30)*time.Second))

		//前10秒已封盘
		qsLockedQueue, _ := NewQsLockedTask(&qs)
		MClient.Enqueue(qsLockedQueue, asynq.ProcessIn(time.Duration(qsTime-20)*time.Second))

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
		if newQs.TaskId == "" {
			qsStartQueue, _ := NewQsStartTask(newQs)
			taskinfo, err := MClient.Enqueue(qsStartQueue, asynq.ProcessIn(time.Duration(2)*time.Second))
			newQs.TaskId = taskinfo.ID
			model.EditQs(newQs)
			if err != nil {
				log.Sugar.Error(err)
			}
		}

		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()

	if qs.ID > 0 && qs.Status == 2 {
		//下注时间
		openFixGameStr := "<b>第<code>%s</code>期-开奖时间：%s\n\n—— —— ——封盘线—— —— —— \n\n</b>"
		list, err := model.GetOrderByQsIdAndStatus(qs.ID, 0)
		if err != nil {
			return err
		}
		endTime := time.Unix(qs.EndTime, 0)
		xzGameStrTemp := fmt.Sprintf(openFixGameStr, qs.Sn, endTime.Format("13:04:05"))

		if len(*list) > 0 {
			xzGameStrTemp += "投注玩家\n"
			for _, item := range *list {
				if item.Stake == 1 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "大", item.Money, item.Rate) + "\n"
				} else if item.Stake == 2 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "小", item.Money, item.Rate) + "\n"
				} else if item.Stake >= 3 && item.Stake <= 18 {
					xzGameStrTemp += fmt.Sprintf("%s  %d杀 %.2f  （赔率 1:%.2f）", item.Nickname, item.Stake, item.Money, item.Rate) + "\n"
				} else if item.Stake == 19 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "单", item.Money, item.Rate) + "\n"
				} else if item.Stake == 20 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "双", item.Money, item.Rate) + "\n"
				} else if item.Stake == 21 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "大单", item.Money, item.Rate) + "\n"
				} else if item.Stake == 22 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "大双", item.Money, item.Rate) + "\n"
				} else if item.Stake == 23 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "小单", item.Money, item.Rate) + "\n"
				} else if item.Stake == 24 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "小双", item.Money, item.Rate) + "\n"
				} else if item.Stake == 25 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "对子", item.Money, item.Rate) + "\n"
				} else if item.Stake == 26 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "顺子", item.Money, item.Rate) + "\n"
				} else if item.Stake == 27 {
					xzGameStrTemp += fmt.Sprintf("%s  %s %.2f  （赔率 1:%.2f）", item.Nickname, "豹子", item.Money, item.Rate) + "\n"
				}
			}
			xzGameStrTemp += "\n—已封盘，线上下注全部有效—"
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
		qs.Dx = 1
		qs.Ds = 1
		qs.Dz = 2
		qs.Sz = 2
		qs.Bz = 2
		if diceSum < 11 {
			diceDx = "小"
			qs.Dx = 2
		}
		if diceSum%2 == 0 {
			diceDs = "双"
			qs.Ds = 2
		}

		zjTemp := ""
		//处理中奖结果
		if len(*list) > 0 {
			tx := dao.Mdb.Begin()
			xzType := ""
			for _, item := range *list {
				item.Status = 2
				item.ResMoney = item.Money
				item.Dx = 1
				item.Ds = 1
				item.Dz = 2
				item.Sz = 2
				item.Bz = 2
				if diceDx == "小" {
					item.Dx = 2
				}
				if diceDs == "双" {
					item.Ds = 2
				}
				if dice1Val == dice2Val || dice1Val == dice3Val || dice2Val == dice3Val {
					item.Dz = 1
				}
				if dice1Val == dice2Val && dice1Val == dice3Val {
					item.Bz = 1
				}
				diceArr := []int{dice1Val, dice2Val, dice3Val}
				sort.Ints(diceArr)
				if diceArr[0]+1 == diceArr[1] && diceArr[1]+1 == diceArr[2] {
					item.Sz = 1
				}

				if item.Stake == 1 && item.Dx == 1 {
					item.Status = 1
					xzType = "大"
				} else if item.Stake == 2 && item.Dx == 2 {
					item.Status = 1
					xzType = "小"
				} else if item.Stake >= 3 && item.Stake <= 18 && diceSum == item.Stake {
					item.Status = 1
					item.Sum = diceSum
					xzType = fmt.Sprintf("%d杀", diceSum)
				} else if item.Stake == 19 && item.Ds == 1 {
					item.Status = 1
					xzType = "单"
				} else if item.Stake == 20 && item.Ds == 2 {
					item.Status = 1
					xzType = "双"
				} else if item.Stake == 21 && item.Dx == 1 && item.Ds == 1 {
					item.Status = 1
					xzType = "大单"
				} else if item.Stake == 22 && item.Dx == 1 && item.Ds == 2 {
					item.Status = 1
					xzType = "大双"
				} else if item.Stake == 23 && item.Dx == 2 && item.Ds == 1 {
					item.Status = 1
					xzType = "小单"
				} else if item.Stake == 24 && item.Dx == 2 && item.Ds == 2 {
					item.Status = 1
					xzType = "小双"
				} else if item.Stake == 27 && item.Bz == 1 {
					item.Status = 1
					xzType = "豹子"
				} else if item.Stake == 25 && item.Dz == 1 {
					item.Status = 1
					xzType = "对子"
				} else if item.Stake == 26 && item.Sz == 1 {
					item.Status = 1
					xzType = "顺子"
				}

				item.Res = fmt.Sprintf("%d,%d,%d", dice1Val, dice2Val, dice3Val)
				item.Sum = diceSum
				if item.Status == 1 {
					item.ResMoney = item.Money*item.Rate + item.Money
					user, err := model.GetUserInfoById(item.UserId)
					if err != nil {
						tx.Rollback()
						return err
					}
					user.Money = user.Money + item.ResMoney
					user.FreezMoney = user.FreezMoney - item.Money
					_, err = model.EditUser(user)
					if err != nil {
						tx.Rollback()
						return err
					}
					rateStr := fmt.Sprintf("1:%.2f", item.Rate)
					billData := &model.Bill{
						UserId:   item.UserId,
						TgId:     item.TgId,
						Username: item.Username,
						Nickname: item.Nickname,
						Type:     4,
						ResId:    item.ID,
						Money:    item.ResMoney,
						Remark:   fmt.Sprintf("%s[第%s期中奖]%s - %.2f （赔率 %s）中奖金额：%.2f", item.Nickname, item.QsSn, xzType, item.Money, rateStr, item.ResMoney),
					}
					_, err = model.AddBill(billData)
					if err != nil {
						tx.Rollback()
						return err
					}
					zjTemp += fmt.Sprintf("%s %s - %.2f （赔率 %s）中奖金额：%.2f\n", user.Nickname, xzType, item.Money, rateStr, item.ResMoney)
				} else {
					user, err := model.GetUserInfoById(item.ID)
					if err != nil {
						tx.Rollback()
						return err
					}
					user.FreezMoney = user.FreezMoney - item.Money
					_, err = model.EditUser(user)
					if err != nil {
						tx.Rollback()
						return err
					}
				}
				_, err := model.EditOrder(&item)
				if err != nil {
					tx.Rollback()
					return err
				}
				tx.Commit()
			}
			xzGameStrTemp += "\n—已封盘，线上下注全部有效—"
		}
		//开奖通知
		openGameStr := "<b>第<code>%s</code>期开奖结果：\n%s \n\n🎉🎉恭喜以下中奖玩家🎉🎉\n</b>%s"
		resGames := fmt.Sprintf("%d %d %d = %d %s %s", dice1Val, dice2Val, dice3Val, diceSum, diceDx, diceDs)
		telegram.SendToBotInBtns(fmt.Sprintf(openGameStr, qs.Sn, resGames, zjTemp), qs.Sn)

		//最近10期结果通知
		qsListStr := "Telegram 官方骰子，具体玩法看置顶\n\n----最近10期结果----\n"
		qsList, _ := model.GetQsListByStatus(3, 10)
		for _, item := range qsList {
			kjType := ""
			if item.Dx == 1 {
				kjType += "大"
			} else {
				kjType += "小"
			}
			if item.Ds == 1 {
				kjType += " 单"
			} else {
				kjType += " 双"
			}
			if item.Dz == 1 {
				kjType += " 对子"
			}
			if item.Sz == 1 {
				kjType += " 顺子"
			}
			if item.Bz == 1 {
				kjType += " 豹子"
			}
			qsListStr += fmt.Sprintf("<code>%s</code>期 %s = %d %s \n", item.Sn, item.Res, item.Sum, kjType)
		}
		telegram.SendToBotInBtns(qsListStr, qs.Sn)

		//更新期数
		qs.Status = 3
		qs.Sum = diceSum
		diceArr := []int{dice1Val, dice2Val, dice3Val}
		sort.Ints(diceArr)
		if dice1Val == dice2Val && dice1Val == dice3Val {
			qs.Dz = 1
		} else if dice1Val == dice2Val || dice1Val == dice3Val || dice2Val == dice3Val {
			qs.Dz = 1
		} else if diceArr[0]+1 == diceArr[1] && diceArr[1]+1 == diceArr[2] {
			qs.Sz = 1
		}
		qs.Res = fmt.Sprintf("%d,%d,%d", dice1Val, dice2Val, dice3Val)
		model.EditQs(&qs)

	} else {
		return errors.New("结算失败")
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
	mutexLock.Lock()
	var qs model.Qs
	err := json.Unmarshal(t.Payload(), &qs)
	if err != nil {
		return err
	}
	defer func() {
		mutexLock.Unlock()
		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()
	if qs.ID > 0 && qs.Status == 1 {
		//游戏已封盘，请勿投注，投注将视为无效投注！
		fpOkGameStr := "<b>第<code>%s</code>期已封盘，请勿投注，投注将视为无效投注！</b>"
		telegram.SendToBot(fmt.Sprintf(fpOkGameStr, qs.Sn))
		//更改状态
		qs.Status = 2
		model.EditQs(&qs)

		//设置缓存
		payload, _ := json.Marshal(qs)
		dao.Rdb.Set(ctx, constant.CacheQsNow, payload, time.Duration(8)*time.Second)

		//已结束
		qsClosedQueue, _ := NewQsClosedTask(&qs)
		MClient.Enqueue(qsClosedQueue, asynq.ProcessIn(time.Duration(8)*time.Second))
	}
	return nil
}
