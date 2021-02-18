package model

import (
	"github.com/jinzhu/gorm"
)

type Problem struct {
	gorm.Model
	GroupId      uint `gorm:"uniqueIndex"`
	RawProblemId uint
	RawProblem   *RawProblem
	Status       int32
	Submitted    int64 `gorm:"default:0"`
	Accepted     int64 `gorm:"default:0"`
}
