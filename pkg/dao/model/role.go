package model

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	RoleName string
}
