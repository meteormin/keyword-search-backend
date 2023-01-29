package entity

import (
	"github.com/miniyus/keyword-search-backend/internal/group_detail"
	"gorm.io/gorm"
	"time"
)

type UserRole string

const (
	Admin UserRole = "admin"
)

type User struct {
	gorm.Model
	Username        string     `gorm:"column:username;type:varchar(50);uniqueIndex" json:"username"`
	Email           string     `gorm:"column:email;type:varchar(100);uniqueIndex" json:"email"`
	Password        string     `gorm:"column:password;type:varchar(255)" json:"-"`
	GroupId         *uint      `gorm:"column:group_id;type:bigint" json:"group_id"`
	Role            UserRole   `gorm:"column:role;type:varchar(10)" json:"role"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at" json:"email_verified_at"`
	Group           Group
	Hosts           []*Host     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"hosts"`
	BookMarks       []*BookMark `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"bookmarks"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	groupDetailHandler := group_detail.HandleCreatedUser(group_detail.CreateParameter{
		DB: tx,
	})

	return groupDetailHandler(u, tx)
}
