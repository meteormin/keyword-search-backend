package entity

import (
	dataTypes "github.com/miniyus/keyword-search-backend/internal/entity/dataTypes"
	"gorm.io/gorm"
)

type Action struct {
	action string
	table  string
}

type Permission struct {
	gorm.Model
	Group      *Group         `json:"group" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	GroupId    uint           `json:"group_id"`
	Permission string         `json:"permission" gorm:"column:permission;type:varchar(10)"`
	Action     dataTypes.JSON `json:"action" gorm:"column:action;type:json"`
}
