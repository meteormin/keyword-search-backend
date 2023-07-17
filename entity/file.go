package entity

import "gorm.io/gorm"

type File struct {
	gorm.Model
	MimeType  string
	Size      int64
	Path      string `gorm:"column:path"`
	Extension string `gorm:"column:extension;"`
}
