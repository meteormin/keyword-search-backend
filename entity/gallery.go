package entity

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	UserId      uint
	User        User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Subject     string  `gorm:"column:subject;type:varchar"`
	Description string  `gorm:"column:description;type:varchar"`
	Photo       []Photo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Photo struct {
	gorm.Model
	Caption   string  `gorm:"column:caption;type:varchar"`
	Gallery   Gallery `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	GalleryId uint
	File      File `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FileId    uint
}
