package entity

import "gorm.io/gorm"

type File struct {
	gorm.Model
	MimeType  string
	Path      string `gorm:"column:path"`
	Extension string `gorm:"column:extension;"`
}
