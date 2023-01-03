package entity

import (
	"gorm.io/gorm"
)

type Action struct {
	gorm.Model
	Method       string `json:"method" gorm:"column:method;type:varchar(10)"`
	Resource     string `json:"resource" gorm:"column:resource;type:varchar(50)"`
	PermissionId uint   `json:"permission_id" gorm:"column:permission_id;type:bigint"`
}

type Permission struct {
	gorm.Model
	Group      *Group   `json:"group" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	GroupId    uint     `json:"group_id"`
	Permission string   `json:"permission" gorm:"column:permission;type:varchar(10)"`
	Actions    []Action `json:"actions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
