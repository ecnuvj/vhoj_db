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

func (c *ContestMapperImpl) BatchSave(contestId uint, problemIds []uint) error {
	var buffer bytes.Buffer
	sql := "insert into `contest_problems` (`contest_id`,`problem_id`) values"
	buffer.WriteString(sql)
	for i, p := range problemIds {
		buffer.WriteString(fmt.Sprintf("(%v,%v)", contestId, p))
		if i == len(problemIds)-1 {
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
	err := c.BatchSave(contest.ID, contest.ProblemIds)
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
		c.DB.Debug().Table("contest_problems").Select("problem_id").Where("contest_id = ?", con.ID).Find(&contestProblems)
		//fmt.Println(contestProblems)
		problemIds := make([]uint, len(contestProblems))
		for ii, cp := range contestProblems {
			problemIds[ii] = cp.ProblemId
		}
		contests[i].ProblemIds = problemIds
	}
	return contests, count, nil
}
