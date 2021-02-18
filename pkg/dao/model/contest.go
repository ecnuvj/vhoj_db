package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Contest struct {
	gorm.Model
	Title       string
	Description string
	UserId      uint
	User        *User
	ProblemNum  int64
	ProblemIds  []uint `gorm:"-"`
	StartTime   time.Time
	EndTime     time.Time
}
