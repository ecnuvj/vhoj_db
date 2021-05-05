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
	AddContestProblemSubmittedCountById(contestId uint, problemId uint) error
	AddContestProblemAcceptedCountById(contestId uint, problemId uint) error
	AddOrModifyProblemGroup(*model.ProblemGroup) (*model.ProblemGroup, error)
	AddOrModifyProblem(*model.Problem) (*model.Problem, error)
	UpdateProblemGroupId(uint, uint) error
	FindGroupProblemsById(uint) ([]*model.ProblemGroup, error)
	FindAllProblems(int32, int32, bool) ([]*model.Problem, int32, error)
	FindProblemById(uint) (*model.Problem, error)
	FindProblemsByIds([]uint) ([]*model.Problem, error)
	SearchProblemByCondition(*ProblemSearchParam, int32, int32) ([]*model.Problem, int32, error)
	DeleteProblemById(uint) error
	FindProblemByRandom() (*model.Problem, error)
	FindRawProblemsWithGroup(int32, int32) ([]*model.RawProblem, []*model.ProblemGroup, int32, error)
	UpdateProblemGroup(uint, uint) error
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
			Model(rawProblem).
			Where("remote_oj = ? and remote_problem_id = ?", rawProblem.RemoteOJ, rawProblem.RemoteProblemId).
			Update(rawProblem).
			First(rawProblem)
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

func (p *ProblemMapperImpl) FindAllProblems(page int32, pageSize int32, desc bool) ([]*model.Problem, int32, error) {
	limit, offset := util.CalLimitOffset(page, pageSize)
	var count int32
	var problems []*model.Problem
	result := p.DB.
		Model(&model.Problem{}).
		Count(&count).
		Preload("RawProblem").
		Limit(limit).
		Offset(offset)
	if desc {
		result = result.Order("updated_at desc")
	}
	result.Find(&problems)
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
		return p.FindAllProblems(pageNo, pageSize, false)
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

func (p *ProblemMapperImpl) AddOrModifyProblemGroup(group *model.ProblemGroup) (*model.ProblemGroup, error) {
	var tmpGroup model.ProblemGroup
	if err := p.DB.Where("raw_problem_id = ?", group.RawProblemId).First(&tmpGroup).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			result := p.DB.Create(group)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	} else {
		result := p.DB.
			Model(group).
			Where("raw_problem_id = ?", group.RawProblemId).
			Update(group).
			First(group)
		if result.Error != nil {
			return nil, result.Error
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

func (p *ProblemMapperImpl) AddOrModifyProblem(problem *model.Problem) (*model.Problem, error) {
	var tmpProblem model.Problem
	if err := p.DB.Where("group_id = ?", problem.GroupId).First(&tmpProblem).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			result := p.DB.Create(problem)
			if result.Error != nil {
				return nil, result.Error
			}
		}
	} else {
		result := p.DB.
			Model(problem).
			Where("group_id = ?", problem.GroupId).
			Update(problem).
			First(problem)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return problem, nil
}

func (p *ProblemMapperImpl) AddContestProblemSubmittedCountById(contestId uint, problemId uint) error {
	result := p.DB.
		Table("contest_problems").
		Where("contest_id = ? and problem_id = ?", contestId, problemId).
		Update("submitted", gorm.Expr("submitted + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *ProblemMapperImpl) AddContestProblemAcceptedCountById(contestId uint, problemId uint) error {
	result := p.DB.
		Table("contest_problems").
		Where("contest_id = ? and problem_id = ?", contestId, problemId).
		Update("accepted", gorm.Expr("accepted + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *ProblemMapperImpl) DeleteProblemById(problemId uint) error {
	problem := &model.Problem{
		Model: gorm.Model{
			ID: problemId,
		},
	}
	return p.DB.Delete(problem).Error
}

func (p *ProblemMapperImpl) FindProblemByRandom() (*model.Problem, error) {
	var problem model.Problem
	result := p.DB.
		Model(&problem).
		Order("rand()").
		Limit(1).
		Find(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &problem, nil
}

func (p *ProblemMapperImpl) FindRawProblemsWithGroup(pageNo int32, pageSize int32) ([]*model.RawProblem, []*model.ProblemGroup, int32, error) {
	limit, offset := util.CalLimitOffset(pageNo, pageSize)
	var rawProblems []*model.RawProblem
	var count int32
	err := p.DB.Debug().
		Model(&model.RawProblem{}).
		Count(&count).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Find(&rawProblems).
		Error
	if err != nil {
		return nil, nil, 0, err
	}
	rawProblemIds := make([]uint, 0, len(rawProblems))
	for _, rawProblem := range rawProblems {
		rawProblemIds = append(rawProblemIds, rawProblem.ID)
	}
	var problemGroups []*model.ProblemGroup
	err = p.DB.
		Model(&model.ProblemGroup{}).
		Where("raw_problem_id in (?)", rawProblemIds).
		Find(&problemGroups).Error
	if err != nil {
		return nil, nil, 0, err
	}
	return rawProblems, problemGroups, count, nil
}
func (p *ProblemMapperImpl) UpdateProblemGroup(rawProblemId uint, groupId uint) error {
	tx := p.DB.Debug().Begin()
	err := tx.
		Model(&model.ProblemGroup{}).
		Where("raw_problem_id = ?", rawProblemId).
		Update("group_id", groupId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var problem model.Problem
	if err = tx.Where("group_id = ?", groupId).First(&problem).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			problem = model.Problem{
				GroupId:      groupId,
				RawProblemId: rawProblemId,
			}
			if err = tx.Create(&problem).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
