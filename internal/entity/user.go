package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username        string    `gorm:"column:username;type:varchar(50);uniqueIndex"`
	Email           string    `gorm:"column:email;type:varchar(100);uniqueIndex"`
	Password        string    `gorm:"column:password;type:varchar(255)"`
	GroupId         uint      `gorm:"column:group_id;type:bigint"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at"`
}
