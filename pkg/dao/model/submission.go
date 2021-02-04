package model

import (
	"github.com/bqxtt/vhoj_common/pkg/common/constants/language"
	"github.com/bqxtt/vhoj_common/pkg/common/constants/remote_oj"
	"github.com/bqxtt/vhoj_common/pkg/common/constants/status_type"
	"github.com/jinzhu/gorm"
)

type Submission struct {
	gorm.Model
	SubmissionCode *SubmissionCode
	ProblemId      uint
	UserId         uint
	Result         status_type.SubmissionStatusType
	TimeCost       int64
	MemoryCost     int64
	Language       language.Language
	ContestId      uint
	SourceOj       remote_oj.RemoteOJ
	RealRunId      string
}
