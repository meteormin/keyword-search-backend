package entity

import (
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Host        string `gorm:"column:host;type:varchar(100);uniqueIndex" json:"host"`
	Subject     string `gorm:"column:subject;type:varchar(100)" json:"subject"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	Path        string `gorm:"column:description;type:varchar(255)" json:"path"`
	Publish     bool   `gorm:"column:publish;type:bool" json:"publish"`
}
