package model

import "github.com/ecnuvj/vhoj_db/pkg/dao/model/base"

type RoleAuthority struct {
	base.Model
	RoleId      uint
	AuthorityId uint
}
