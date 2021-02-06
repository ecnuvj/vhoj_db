package model

type CompileInfo struct {
	SubmissionId uint
	Info         string `gorm:"type:text"`
}
