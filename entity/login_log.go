package entity

import (
	"gorm.io/gorm"
	"time"
)

type LoginLog struct {
	gorm.Model
	UserId  uint      `json:"user_id"`
	User    User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user"`
	LoginAt time.Time `gorm:"column:login_at" json:"login_at"`
	Ip      string    `gorm:"column:ip" json:"ip"`
}
