package model

import (
	"github.com/ecnuvj/vhoj_db/pkg/dao/model/base"
	"time"
)

type Contest struct {
	base.Model
	ContestId   uint `gorm:"primary_key"`
	Title       string
	Description string
	UserId      uint
	ProblemNum  int64
	StartTime   time.Time
	EndTime     time.Time
}
