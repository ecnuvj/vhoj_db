package model

type ContestProblem struct {
	ContestId    uint   `gorm:"unique_index:uni_idx_pod"`
	ProblemOrder string `gorm:"unique_index:uni_idx_pod"`
	ProblemId    uint
	Title        string
	Submitted    uint `gorm:"default:0"`
	Accepted     uint `gorm:"default:0"`
}
