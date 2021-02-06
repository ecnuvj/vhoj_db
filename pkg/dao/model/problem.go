package model

import (
	"github.com/jinzhu/gorm"
)

type Problem struct {
	gorm.Model
	GroupId      uint
	RawProblemId uint
	Status       int32
	Submitted    int64
	Accepted     int64
}
