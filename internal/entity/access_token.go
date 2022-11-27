package entity

import (
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	gorm.Model
	UserId    uint      `gorm:"column:user_id;foreignKey;references:users"`
	Token     string    `gorm:"column:token;type:varchar(255);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null;index"`
}
