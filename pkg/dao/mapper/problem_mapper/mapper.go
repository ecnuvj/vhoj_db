package problem_mapper

import (
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/ecnuvj/vhoj_db/pkg/util"
	"github.com/jinzhu/gorm"
)

type IProblemMapper interface {
	AddOrModifyRawProblem(*model.RawProblem) (*model.RawProblem, error)
	FindGroupProblemsById(uint) ([]*model.ProblemGroup, error)
	FindAllProblems(page int32, pageSize int32) ([]*model.Problem, uint32, error)
}

var ProblemMapper IProblemMapper

type ProblemMapperImpl struct {
	DB *gorm.DB
}

func InitMapper(db *gorm.DB) {
	ProblemMapper = &ProblemMapperImpl{
		DB: db,
	}
}

func (p *ProblemMapperImpl) AddOrModifyRawProblem(rawProblem *model.RawProblem) (*model.RawProblem, error) {
	var problem model.RawProblem
	if err := p.DB.Where("remote_oj = ? and remote_problem_id = ?", rawProblem.RemoteOJ, rawProblem.RemoteProblemId).First(&problem).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			result := p.DB.Create(rawProblem)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	} else {
		result := p.DB.
			Model(&rawProblem).
			Where("remote_oj = ? and remote_problem_id = ?", rawProblem.RemoteOJ, rawProblem.RemoteProblemId).
			Update(rawProblem)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return rawProblem, nil
}

func (p *ProblemMapperImpl) FindGroupProblemsById(problemId uint) ([]*model.ProblemGroup, error) {
	var problemGroups []*model.ProblemGroup
	var problem model.Problem
	result := p.DB.
		Select("group_id").
		Where("id = ?", problemId).
		First(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	result = p.DB.
		Model(&model.ProblemGroup{}).
		Where("group_id = ?", problem.GroupId).
		Find(&problemGroups)
	if result.Error != nil {
		return nil, result.Error
	}
	return problemGroups, nil
}

func (p *ProblemMapperImpl) FindAllProblems(page int32, pageSize int32) ([]*model.Problem, uint32, error) {
	limit, offset := util.CalLimitOffset(page, pageSize)
	var count uint32
	var problems []*model.Problem
	result := p.DB.
		Model(&model.Problem{}).
		Count(&count).
		Preload("RawProblem").
		Limit(limit).
		Offset(offset).
		Find(&problems)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return problems, count, nil
}
