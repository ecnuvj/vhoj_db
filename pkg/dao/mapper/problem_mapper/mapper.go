package problem_mapper

import (
	"fmt"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/ecnuvj/vhoj_db/pkg/util"
	"github.com/jinzhu/gorm"
)

type ProblemSearchParam struct {
	Title     string
	ProblemId uint
}

type IProblemMapper interface {
	AddOrModifyRawProblem(*model.RawProblem) (*model.RawProblem, error)
	AddProblemSubmittedCountById(uint) error
	AddProblemAcceptedCountById(uint) error
	AddProblemGroup(*model.ProblemGroup) (*model.ProblemGroup, error)
	UpdateProblemGroupId(uint, uint) error
	FindGroupProblemsById(uint) ([]*model.ProblemGroup, error)
	FindAllProblems(int32, int32) ([]*model.Problem, int32, error)
	FindProblemById(uint) (*model.Problem, error)
	FindProblemsByIds([]uint) ([]*model.Problem, error)
	SearchProblemByCondition(*ProblemSearchParam, int32, int32) ([]*model.Problem, int32, error)
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

func (p *ProblemMapperImpl) FindAllProblems(page int32, pageSize int32) ([]*model.Problem, int32, error) {
	limit, offset := util.CalLimitOffset(page, pageSize)
	var count int32
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

func (p *ProblemMapperImpl) AddProblemSubmittedCountById(problemId uint) error {
	if problemId <= 0 {
		return fmt.Errorf("problem id is incorrect")
	}
	result := p.DB.
		Model(&model.Problem{Model: gorm.Model{ID: problemId}}).
		Update("submitted", gorm.Expr("submitted + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (p *ProblemMapperImpl) AddProblemAcceptedCountById(problemId uint) error {
	if problemId <= 0 {
		return fmt.Errorf("problem id is incorrect")
	}
	result := p.DB.
		Model(&model.Problem{Model: gorm.Model{ID: problemId}}).
		Update("accepted", gorm.Expr("accepted + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *ProblemMapperImpl) FindProblemById(problemId uint) (*model.Problem, error) {
	if problemId <= 0 {
		return nil, fmt.Errorf("problem id is incorrect")
	}
	var problem = &model.Problem{
		Model: gorm.Model{
			ID: problemId,
		},
		RawProblem: &model.RawProblem{},
	}
	result := p.DB.
		Model(problem).
		First(problem).
		Related(problem.RawProblem)
	if result.Error != nil {
		return nil, result.Error
	}
	return problem, nil
}

func (p *ProblemMapperImpl) FindProblemsByIds(problemIds []uint) ([]*model.Problem, error) {
	var problems []*model.Problem
	result := p.DB.
		Model(&model.Problem{}).
		Preload("RawProblem").
		Find(&problems, problemIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return problems, nil
}

func (p *ProblemMapperImpl) SearchProblemByCondition(param *ProblemSearchParam, pageNo int32, pageSize int32) ([]*model.Problem, int32, error) {
	if (param.ProblemId == 0 && param.Title == "") || param == nil {
		return p.FindAllProblems(pageNo, pageSize)
	}
	limit, offset := util.CalLimitOffset(pageNo, pageSize)
	result := p.DB
	if param.Title != "" {
		result = result.Preload("RawProblem", "title LIKE ?", fmt.Sprintf("%%%v%%", param.Title))
	} else {
		result = result.Preload("RawProblem")
	}
	if param.ProblemId != 0 {
		result = result.Where("id = ?", param.ProblemId)
	}
	var problems []*model.Problem
	result = result.Find(&problems)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	retProblems := make([]*model.Problem, 0)
	for _, problem := range problems {
		if problem.RawProblem != nil {
			retProblems = append(retProblems, problem)
		}
	}
	left, right := util.CalSliceLeftRight(limit, offset, int32(cap(retProblems)))
	return retProblems[left:right], int32(len(retProblems)), nil
}

func (p *ProblemMapperImpl) AddProblemGroup(group *model.ProblemGroup) (*model.ProblemGroup, error) {
	if err := p.DB.Where("raw_problem_id = ?", group.RawProblemId).Find(group).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			p.DB.Create(group)
		}
	}
	return group, nil
}

func (p *ProblemMapperImpl) UpdateProblemGroupId(rawProblemId uint, groupId uint) error {
	result := p.DB.
		Model(&model.ProblemGroup{}).
		Where("raw_problem_id = ?", rawProblemId).
		Update("group_id", groupId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
