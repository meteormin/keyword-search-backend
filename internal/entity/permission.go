package entity

import (
	dataTypes "github.com/miniyus/go-fiber/internal/entity/dataTypes"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Permission string         `json:"permission" gorm:"column:permission;type:varchar(10)"`
	Action     dataTypes.JSON `json:"action" gorm:"column:action;type:json"`
}
