package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type RuntimeInfo struct {
	base.Model
	SubmissionId uint
	Info         string
}
