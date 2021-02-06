package model

import (
	"github.com/jinzhu/gorm"
)

type ProblemGroup struct {
	gorm.Model
	RawProblemId uint
	GroupId      uint
	MainProblem  bool
}
