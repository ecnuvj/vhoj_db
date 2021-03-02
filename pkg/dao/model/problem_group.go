package model

import (
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/remote_oj"
	"github.com/jinzhu/gorm"
)

type ProblemGroup struct {
	gorm.Model
	RawProblemId    uint `gorm:"unique_index:uidx_rawId"`
	GroupId         uint
	MainProblem     bool
	RemoteOJ        remote_oj.RemoteOJ
	RemoteProblemId string
}
