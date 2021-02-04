package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type ContestProblem struct {
	base.Model
	ContestId uint
	ProblemId uint
}
