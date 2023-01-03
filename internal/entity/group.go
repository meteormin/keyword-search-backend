package entity

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string       `json:"name" gorm:"column:name;type:varchar(50);uniqueIndex"`
	Permissions []Permission `json:"permissions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Users       []User       `json:"users" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
