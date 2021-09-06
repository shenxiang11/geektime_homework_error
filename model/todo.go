package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Content string `gorm:"type:varchar(50);not null"`
	Status int `gorm:"type:tinyint comment '1未完成 2已完成';not null;default:1"`
}