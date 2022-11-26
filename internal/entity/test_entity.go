package entity

import "gorm.io/gorm"

type TestEntity struct {
	gorm.Model
	TestField string
}
