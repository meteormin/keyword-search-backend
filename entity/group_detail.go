package entity

import "gorm.io/gorm"

type GroupRole string

const (
	Owner   GroupRole = "owner"
	Manager GroupRole = "manager"
	Member  GroupRole = "member"
)

type GroupDetail struct {
	gorm.Model
	GroupId uint `json:"group_id" gorm:"column:group_id;type:bigint"`
	Group   Group
	UserId  uint `json:"user_id" gorm:"column:user_id;type:bigint"`
	User    User
	Role    GroupRole `json:"role" gorm:"column:role;type:varchar(10)"`
}
