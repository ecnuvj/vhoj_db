package model

import (
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/remote_oj"
	"github.com/jinzhu/gorm"
)

type RawProblem struct {
	gorm.Model
	Title           string
	Description     string             `gorm:"type:text"`
	SampleInput     string             `gorm:"type:text"`
	SampleOutput    string             `gorm:"type:text"`
	Input           string             `gorm:"type:text"`
	Output          string             `gorm:"type:text"`
	Hint            string             `gorm:"type:text"`
	RemoteOJ        remote_oj.RemoteOJ `gorm:"unique_index:uni_idx_pid"`
	RemoteProblemId string             `gorm:"unique_index:uni_idx_pid"`
	TimeLimit       string
	MemoryLimit     string
	Spj             string
	Std             string
	Source          string
}
