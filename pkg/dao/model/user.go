package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserAuth     *UserAuth
	Email        string
	Nickname     string `gorm:"unique_index:uidx_name"`
	School       string
	Roles        []*Role `gorm:"-"`
	GenerateUser bool
	ContestId    uint
	Submitted    int64
	Accepted     int64
}
