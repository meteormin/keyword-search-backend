package entity

import "gorm.io/gorm"

type Search struct {
	gorm.Model
	Host        *Host   `json:"host" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HostId      uint    `json:"host_id" gorm:"index:query_unique,unique"`
	QueryKey    string  `gorm:"column:query_key;type:varchar(50);index:query_unique,unique" json:"query_key"`
	Query       string  `gorm:"column:query;type:varchar(50);index:query_unique,unique" json:"query"`
	Description string  `gorm:"column:description;type:varchar(50)" json:"description"`
	Publish     bool    `gorm:"column:publish;type:bool" json:"publish"`
	ShortUrl    *string `gorm:"column:short_url;type:varchar(255);uniqueIndex" json:"short_url"`
}
