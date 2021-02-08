package model

import "github.com/jinzhu/gorm"

type UserAuth struct {
	gorm.Model
	UserID   uint
	Password string
}
