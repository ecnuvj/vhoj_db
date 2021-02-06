package model

type RuntimeInfo struct {
	SubmissionId uint
	Info         string `gorm:"type:text"`
}
