package model

import (
	"ant/utils/dao"
	"time"
)

type Order struct {
	UserId   uint64  `gorm:"column:user_id" json:"user_id"`
	TgId     string  `gorm:"column:tg_id" json:"tg_id"`
	Username string  `gorm:"column:username" json:"username"`
	Nickname string  `gorm:"column:nickname" json:"nickname"`
	QsId     uint64  `gorm:"column:qs_id" json:"qs_id"`
	QsSn     string  `gorm:"column:qs_sn" json:"qs_sn"`
	Rate     float64 `gorm:"column:rate" json:"rate"`
	Stake    int     `gorm:"column:stake" json:"stake"`
	Money    float64 `gorm:"column:money" json:"money"`
	ResMoney float64 `gorm:"column:res_money" json:"res_money"`
	Status   int     `gorm:"column:status" json:"status"`
	Res      string  `gorm:"column:res" json:"res"`
	Sum      int     `gorm:"column:sum" json:"sum"`
	Dx       int     `gorm:"column:dx" json:"dx"`
	Ds       int     `gorm:"column:ds" json:"ds"`
	Dz       int     `gorm:"column:dz" json:"dz"`
	Sz       int     `gorm:"column:sz" json:"sz"`
	Bz       int     `gorm:"column:bz" json:"bz"`
	BaseModel
}

func (order *Order) TableName() string {
	return "fa_order"
}

func OrderList() ([]Order, error) {
	list := []Order{}
	err := dao.Mdb.Model(&Order{}).Find(&list).Error
	return list, err
}

func GetOrderByQsId(qsId uint64) (*[]Order, error) {
	list := &[]Order{}
	err := dao.Mdb.Model(&Order{}).Where("qs_id = ?", qsId).Find(list).Error
	return list, err
}

func GetOrderByQsIdAndStatus(qsId uint64, status int) (*[]Order, error) {
	list := &[]Order{}
	err := dao.Mdb.Model(&Order{}).Where("qs_id = ?", qsId).Where("status = ?", status).Find(list).Error
	return list, err
}

func AddOrder(order *Order) (*Order, error) {
	order.CreateTime = time.Now().Unix()
	err := dao.Mdb.Model(&Order{}).Create(order).Error
	return order, err
}

func EditOrder(order *Order) (*Order, error) {
	order.UpdateTime = time.Now().Unix()
	err := dao.Mdb.Model(order).Where("id=?", order.ID).Updates(order).Error
	return order, err
}
