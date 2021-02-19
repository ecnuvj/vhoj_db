package contest_mapper

import (
	"bytes"
	"fmt"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/ecnuvj/vhoj_db/pkg/util"
	"github.com/jinzhu/gorm"
)

type IContestMapper interface {
	CreateContest(*model.Contest) (*model.Contest, error)
	FindAllContests(int32, int32) ([]*model.Contest, int32, error)
	FindContestById(uint) (*model.Contest, error)
	AddContestParticipants(uint, []uint) error
	AddContestAdmins(uint, []uint) error
	FindContestAdmins(uint) ([]uint, error)
	FindContestParticipants(uint) ([]uint, error)
}

var ContestMapper IContestMapper

type ContestMapperImpl struct {
	DB *gorm.DB
}

func InitMapper(db *gorm.DB) {
	ContestMapper = &ContestMapperImpl{
		DB: db,
	}
}

func (c *ContestMapperImpl) BatchSave(contestId uint, tableName string, column string, Ids []uint) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into `%v` (`contest_id`,`%v`) values", tableName, column)
	buffer.WriteString(sql)
	for i, id := range Ids {
		buffer.WriteString(fmt.Sprintf("(%v,%v)", contestId, id))
		if i == len(Ids)-1 {
			buffer.WriteString(";")
		} else {
			buffer.WriteString(",")
		}
	}
	return c.DB.Exec(buffer.String()).Error
}

func (c *ContestMapperImpl) CreateContest(contest *model.Contest) (*model.Contest, error) {
	result := c.DB.Create(contest)
	if result.Error != nil {
		return nil, result.Error
	}
	err := c.BatchSave(contest.ID, "contest_problems", "problem_id", contest.ProblemIds)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (c *ContestMapperImpl) FindAllContests(pageNo int32, PageSize int32) ([]*model.Contest, int32, error) {
	limit, offset := util.CalLimitOffset(pageNo, PageSize)
	var count int32
	var contests []*model.Contest
	result := c.DB.
		Model(&model.Contest{}).
		Count(&count).
		Preload("User").
		Limit(limit).
		Offset(offset).
		Find(&contests)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	for i, con := range contests {
		var contestProblems []*model.ContestProblem
		c.DB.Table("contest_problems").Select("problem_id").Where("contest_id = ?", con.ID).Find(&contestProblems)
		//fmt.Println(contestProblems)
		problemIds := make([]uint, len(contestProblems))
		for ii, cp := range contestProblems {
			problemIds[ii] = cp.ProblemId
		}
		contests[i].ProblemIds = problemIds
	}
	return contests, count, nil
}

func (c *ContestMapperImpl) FindContestById(contestId uint) (*model.Contest, error) {
	contest := &model.Contest{
		Model: gorm.Model{
			ID: contestId,
		},
	}
	result := c.DB.Model(contest).Preload("User").First(contest)
	if result.Error != nil {
		return nil, result.Error
	}
	var contestProblems []*model.ContestProblem
	c.DB.
		Table("contest_problems").
		Select("problem_id").
		Where("contest_id = ?", contestId).
		Find(&contestProblems)
	problemIds := make([]uint, len(contestProblems))
	for i, c := range contestProblems {
		problemIds[i] = c.ProblemId
	}
	contest.ProblemIds = problemIds
	return contest, nil
}

func (c *ContestMapperImpl) AddContestParticipants(contestId uint, userIds []uint) error {
	err := c.BatchSave(contestId, "contest_participants", "user_id", userIds)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContestMapperImpl) AddContestAdmins(contestId uint, userIds []uint) error {
	err := c.BatchSave(contestId, "contest_admins", "user_id", userIds)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContestMapperImpl) FindContestAdmins(contestId uint) ([]uint, error) {
	var contestAdmins []*model.ContestAdmin
	result := c.DB.
		Table("contest_admins").
		Where("contest_id = ?", contestId).
		Find(&contestAdmins)
	if result.Error != nil {
		return nil, result.Error
	}
	userIds := make([]uint, len(contestAdmins))
	for i, u := range contestAdmins {
		userIds[i] = u.UserId
	}
	return userIds, nil
}

func (c *ContestMapperImpl) FindContestParticipants(contestId uint) ([]uint, error) {
	var contestParticipants []*model.ContestParticipant
	result := c.DB.
		Table("contest_participants").
		Where("contest_id = ?", contestId).
		Find(&contestParticipants)
	if result.Error != nil {
		return nil, result.Error
	}
	userIds := make([]uint, len(contestParticipants))
	for i, u := range contestParticipants {
		userIds[i] = u.UserId
	}
	return userIds, nil
}
