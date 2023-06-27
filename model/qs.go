package model

import (
	"ant/utils/dao"
	"time"
)

type Qs struct {
	Sn        string `gorm:"column:sn" json:"sn"`
	BeginTime int64  `gorm:"column:begin_time" json:"begin_time"`
	EndTime   int64  `gorm:"column:end_time" json:"end_time"`
	Status    int    `gorm:"column:status" json:"status"`
	Res       string `gorm:"column:res" json:"res"`
	Sum       int    `gorm:"column:sum" json:"sum"`
	Dx        int    `gorm:"column:dx" json:"dx"`
	Ds        int    `gorm:"column:ds" json:"ds"`
	Dz        int    `gorm:"column:dz" json:"dz"`
	Sz        int    `gorm:"column:sz" json:"sz"`
	Bz        int    `gorm:"column:bz" json:"bz"`
	TaskId    string `gorm:"column:task_id" json:"task_id"`
	BaseModel
}

func (qs *Qs) TableName() string {
	return "fa_qs"
}

func QsList() ([]Qs, error) {
	list := []Qs{}
	err := dao.Mdb.Model(&Qs{}).Find(&list).Error
	return list, err
}

func GetQsListByStatus(status int, limit int) ([]Qs, error) {
	list := []Qs{}
	err := dao.Mdb.Model(&Qs{}).Where("status = ?", status).Limit(limit).Order("id desc").Find(&list).Error
	return list, err
}

func GetQsListByTime(time int64) (*Qs, error) {
	qs := new(Qs)
	err := dao.Mdb.Model(qs).Where("begin_time <= ?", time).Where("end_time > ?", time).Find(qs).Error
	return qs, err
}

func AddQsInBatches(qslist *[]Qs) error {
	return dao.Mdb.Model(&Qs{}).CreateInBatches(qslist, 10).Error
}

func EditQs(qs *Qs) (*Qs, error) {
	qs.UpdateTime = time.Now().Unix()
	err := dao.Mdb.Model(qs).Where("id=?", qs.ID).Updates(qs).Error
	return qs, err
}
