package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserAuth  *UserAuth
	Email     string
	Nickname  string `gorm:"uniqueIndex"`
	School    string
	Submitted int64
	Accepted  int64
}
