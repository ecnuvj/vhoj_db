package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type Problem struct {
	base.Model
	ProblemId    uint `gorm:"primary_key"`
	GroupId      uint
	RawProblemId uint
	Status       int32
	Submitted    int64
	Accepted     int64
}
