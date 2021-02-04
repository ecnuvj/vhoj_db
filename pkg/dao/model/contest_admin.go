package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type ContestAdmin struct {
	base.Model
	ContestId uint
	UserId    uint
}
