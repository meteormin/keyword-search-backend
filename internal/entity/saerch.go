package entity

import "gorm.io/gorm"

type Search struct {
	gorm.Model
	HostId      uint   `json:"host_id"`
	Host        Host   `json:"host"`
	Path        string `gorm:"column:path;type:varchar(50)" json:"path"`
	Query       string `gorm:"column:path;type:varchar(50)" json:"query"`
	Description string `gorm:"column:path;type:varchar(50)" json:"description"`
	Publish     bool   `gorm:"column:path;type:bool" json:"publish"`
}
