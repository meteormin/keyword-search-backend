package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username        string    `gorm:"column:username;type:varchar(50);uniqueIndex" json:"username"`
	Email           string    `gorm:"column:email;type:varchar(100);uniqueIndex" json:"email"`
	Password        string    `gorm:"column:password;type:varchar(255)" json:"password"`
	GroupId         uint      `gorm:"column:group_id;type:bigint" json:"group_id"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at" json:"email_verified_at"`
}
