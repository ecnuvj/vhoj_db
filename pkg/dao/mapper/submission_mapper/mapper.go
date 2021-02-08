package submission_mapper

import (
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/problem_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/jinzhu/gorm"
)

type ISubmissionMapper interface {
	AddOrModifySubmission(submission *model.Submission) (*model.Submission, error)
	FindSubmissionById(submissionId uint) (*model.Submission, error)
	FindProblemGroupById(submissionId uint) ([]*model.ProblemGroup, error)
	UpdateSubmissionById(submission *model.Submission) (*model.Submission, error)
	UpdateSubmissionCEInfoById(submissionId uint, info string) error
	ResetSubmissionById(submissionId uint) error
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

func (s *SubmissionMapperImpl) AddOrModifySubmission(submission *model.Submission) (*model.Submission, error) {
	var sub model.Submission
	if err := s.DB.Where("id = ?", submission.ID).First(&sub).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			result := s.DB.Create(submission)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	} else {
		result := s.DB.Model(submission).Update(submission)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return submission, nil
}

func (s *SubmissionMapperImpl) UpdateSubmissionById(submission *model.Submission) (*model.Submission, error) {
	result := s.DB.Model(submission).Update(submission).Find(submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return submission, nil
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
		result := s.DB.Model(&compileInfo).Where("submission_id = ?", submissionId).Update(&compileInfo)
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

func (s *SubmissionMapperImpl) FindProblemGroupById(submissionId uint) ([]*model.ProblemGroup, error) {
	submission := model.Submission{
		Model: gorm.Model{
			ID: submissionId,
		},
	}
	result := s.DB.First(submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return problem_mapper.ProblemMapper.FindGroupProblemsById(submission.ProblemId)
}

func (s *SubmissionMapperImpl) ResetSubmissionById(submissionId uint) error {
	//var submission model.Submission
	result := s.DB.
		Model(&model.Submission{
			Model: gorm.Model{
				ID: submissionId,
			},
		}).
		Updates(map[string]interface{}{
			"result":      0,
			"time_cost":   0,
			"memory_cost": 0,
			"remote_oj":   0,
			"real_run_id": "",
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
