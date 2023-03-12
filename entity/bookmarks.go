package entity

import "gorm.io/gorm"

type BookMark struct {
	gorm.Model
	Subject     string `json:"subject" gorm:"column:description;type:varchar(50)"`
	Description string `json:"description" gorm:"column:description;type:varchar(50)"`
	Url         string `json:"url" gorm:"column:url;type:varchar(50)"`
	Publish     bool   `json:"publish" gorm:"column:publish;type:bool"`
	Views       uint   `json:"views" gorm:"column:views;default:0"`
	UserId      uint   `json:"user_id"`
	User        *User  `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
