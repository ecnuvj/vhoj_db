package model

import (
	"github.com/jinzhu/gorm"
)

type RawProblem struct {
	gorm.Model
	Title        string
	Description  string
	SampleInput  string
	SampleOutput string
	Input        string
	Output       string
	Hint         string
	SourceOj     int32
	SourceId     string
	TimeLimit    string
	MemoryLimit  string
	Spj          string
	Std          string
}
