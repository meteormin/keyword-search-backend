package entity

import (
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	gorm.Model
	token     string    `gorm:"column:token;type:varchar(255);uniqueIndex;not null"`
	expiresAt time.Time `gorm:"column:expires_at;not null;index"`
}
