package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type ProblemGroup struct {
	base.Model
	RawProblemId uint
	GroupId      uint
	MainProblem  bool
}
