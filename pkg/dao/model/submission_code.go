package model

import (
	"github.com/jinzhu/gorm"
)

type SubmissionCode struct {
	gorm.Model
	SubmissionID uint
	SourceCode   string `gorm:"type:text"`
	CodeLength   int64
}
