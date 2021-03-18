package submission_mapper

import (
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/language"
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/status_type"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/problem_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/ecnuvj/vhoj_db/pkg/util"
	"github.com/jinzhu/gorm"
)

type SearchSubmissionCondition struct {
	Username  string
	ProblemId uint
	Status    status_type.SubmissionStatusType
	Language  language.Language
}

type UserSubmissionCondition struct {
	UserId    uint
	ProblemId uint
	ContestId uint
}

type ISubmissionMapper interface {
	AddOrModifySubmission(submission *model.Submission) (*model.Submission, error)
	FindSubmissionById(submissionId uint) (*model.Submission, error)
	FindProblemGroupById(submissionId uint) ([]*model.ProblemGroup, error)
	FindSubmissions(pageNo int32, pageSize int32, condition *SearchSubmissionCondition) ([]*model.Submission, int32, error)
	FindSubmissionsGroupByResult(condition *UserSubmissionCondition) ([]*model.Submission, error)
	FindSubmissionsByContestId(uint) ([]*model.Submission, error)
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

func (s *SubmissionMapperImpl) FindSubmissions(pageNo int32, pageSize int32, condition *SearchSubmissionCondition) ([]*model.Submission, int32, error) {
	result := s.DB.Model(&model.Submission{})
	limit, offset := util.CalLimitOffset(pageNo, pageSize)
	if condition == nil {
		condition = &SearchSubmissionCondition{}
	}
	var count int32
	var submissions []*model.Submission
	if condition.Username != "" {
		result = result.Where("username = ?", condition.Username)
	}
	if condition.ProblemId != 0 {
		result = result.Where("problem_id = ?", condition.ProblemId)
	}
	if condition.Status != 0 {
		result = result.Where("result = ?", condition.Status)
	}
	if condition.Language != 0 {
		result = result.Where("language = ?", condition.Language)
	}
	result = result.
		Count(&count).
		Order("updated_at desc").
		Limit(limit).
		Offset(offset).
		Find(&submissions)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return submissions, count, nil
}

func (s *SubmissionMapperImpl) FindSubmissionsGroupByResult(condition *UserSubmissionCondition) ([]*model.Submission, error) {
	var submissions []*model.Submission
	result := s.DB.
		Model(&model.Submission{}).
		Select("result").
		Where("user_id = ? and problem_id = ? and contest_id = ?", condition.UserId, condition.ProblemId, condition.ContestId).
		Group("result").
		Find(&submissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return submissions, nil
}

func (s *SubmissionMapperImpl) FindSubmissionsByContestId(contestId uint) ([]*model.Submission, error) {
	var submission []*model.Submission
	result := s.DB.
		Model(&model.Submission{}).
		Where("contest_id = ?", contestId).
		Find(&submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return submission, nil
}
