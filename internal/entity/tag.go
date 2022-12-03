package entity

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Tag         string `json:"tag" gorm:"column:tag;type:varchar(20)"`
	ParentTable string `json:"parent_table" gorm:"column:parent_table;type:varchar(20)"`
	ParentId    uint   `json:"parent_id" gorm:"column:parent_table"`
}
