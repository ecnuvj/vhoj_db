package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type Authority struct {
	base.Model
	AuthorityId uint `gorm:"primary_key"`
	Privilege   string
}
