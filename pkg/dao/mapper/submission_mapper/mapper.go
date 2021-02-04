package submission_mapper

import (
	"github.com/bqxtt/vhoj_common/pkg/common/constants/status_type"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/jinzhu/gorm"
)

type ISubmissionMapper interface {
	AddSubmission(submission *model.Submission) (*model.Submission, error)
	FindSubmissionById(submissionId uint) (*model.Submission, error)
	UpdateSubmissionResultById(submissionId uint, statusType status_type.SubmissionStatusType) error
	UpdateSubmissionCEInfoById(submissionId uint, info string) error
}

var SubmissionMapper ISubmissionMapper

type SubmissionMapperImpl struct {
	DB *gorm.DB
}

func InitMapper(db *gorm.DB) {
	SubmissionMapper = &SubmissionMapperImpl{
		DB: db,
	}
}

func (s *SubmissionMapperImpl) AddSubmission(submission *model.Submission) (*model.Submission, error) {
	result := s.DB.Create(submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return submission, nil
}

func (s *SubmissionMapperImpl) UpdateSubmissionResultById(submissionId uint, statusType status_type.SubmissionStatusType) error {
	submission := &model.Submission{Model: gorm.Model{ID: submissionId}, Result: statusType}
	result := s.DB.Model(submission).Update(submission)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SubmissionMapperImpl) UpdateSubmissionCEInfoById(submissionId uint, info string) error {
	var compileInfo model.CompileInfo
	compileInfo.SubmissionId = submissionId
	if err := s.DB.Where("submission_id = ?", submissionId).First(&compileInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			compileInfo.Info = info
			result := s.DB.Create(&compileInfo)
			if result.Error != nil {
				return result.Error
			}
		}
	} else {
		result := s.DB.Model(&compileInfo).Update(&compileInfo)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (s *SubmissionMapperImpl) FindSubmissionById(submissionId uint) (*model.Submission, error) {
	submission := &model.Submission{Model: gorm.Model{ID: submissionId}}
	code := &model.SubmissionCode{}
	result := s.DB.Model(submission).Find(submission).Related(&code)
	if result.Error != nil {
		return nil, result.Error
	}
	submission.SubmissionCode = code
	return submission, nil
}
