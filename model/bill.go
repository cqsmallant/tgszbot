package model

import (
	"ant/utils/dao"
	"time"
)

type Bill struct {
	UserId   uint64  `gorm:"column:user_id" json:"user_id"`
	TgId     string  `gorm:"column:tg_id" json:"tg_id"`
	Username string  `gorm:"column:username" json:"username"`
	Nickname string  `gorm:"column:nickname" json:"nickname"`
	Type     int     `gorm:"column:type" json:"type"`
	ResId    uint64  `gorm:"column:res_id" json:"res_id"`
	Money    float64 `gorm:"column:money" json:"money"`
	Remark   string  `gorm:"column:remark" json:"remark"`
	BaseModel
}

func (bill *Bill) TableName() string {
	return "fa_bill"
}

func BillList() ([]Bill, error) {
	list := []Bill{}
	err := dao.Mdb.Model(&Bill{}).Order("id desc").Find(&list).Error
	return list, err
}

// 获取列表
func GetBillByTgId(tgId string, limit int) (*[]Bill, error) {
	list := &[]Bill{}
	err := dao.Mdb.Model(&Bill{}).Where("tg_id = ?", tgId).Limit(limit).Order("id desc").Find(list).Error
	return list, err
}

func GetBillByUserId(userId int64) (*Bill, error) {
	info := new(Bill)
	err := dao.Mdb.Model(info).Where("user_id = ?", userId).Find(info).Error
	return info, err
}

func AddBill(data *Bill) (*Bill, error) {
	data.CreateTime = time.Now().Unix()
	err := dao.Mdb.Model(&Bill{}).Create(data).Error
	return data, err
}

func EditBill(data *Bill) (*Bill, error) {
	data.UpdateTime = time.Now().Unix()
	err := dao.Mdb.Model(data).Where("id=?", data.ID).Updates(data).Error
	return data, err
}
