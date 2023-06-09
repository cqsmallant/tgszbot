package model

type BaseModel struct {
	ID         uint64 `gorm:"column:id;primary_key" json:"id"`
	CreateTime int64  `gorm:"column:createtime" json:"create_time"`
	UpdateTime int64  `gorm:"column:updatetime" json:"update_time"`
}
