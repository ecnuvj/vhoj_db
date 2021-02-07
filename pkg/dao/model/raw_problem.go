package model

import (
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/remote_oj"
	"github.com/jinzhu/gorm"
)

type RawProblem struct {
	gorm.Model
	Title           string
	Description     string
	SampleInput     string
	SampleOutput    string
	Input           string
	Output          string
	Hint            string
	RemoteOJ        remote_oj.RemoteOJ `gorm:"unique_index:uni_idx_pid"`
	RemoteProblemId string             `gorm:"unique_index:uni_idx_pid"`
	TimeLimit       string
	MemoryLimit     string
	Spj             string
	Std             string
}
