package model

type CompileInfo struct {
	SubmissionId uint `gorm:"primary_key"`
	Info         string
}
