package contest_mapper

import (
	"bytes"
	"fmt"
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/contest_status"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/user_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/ecnuvj/vhoj_db/pkg/util"
	"github.com/jinzhu/gorm"
	"time"
)

type SearchContestCondition struct {
	Status      contest_status.ContestStatus
	Title       string
	CreatorName string
}

type IContestMapper interface {
	CreateContest(*model.Contest, []*model.ContestProblem) (*model.Contest, error)
	FindAllContests(int32, int32) ([]*model.Contest, int32, error)
	FindContestById(uint) (*model.Contest, error)
	FindContestsByCondition(*SearchContestCondition, int32, int32) ([]*model.Contest, int32, error)
	FindContestAdmins(uint) ([]uint, error)
	FindContestParticipants(uint) ([]uint, error)
	FindContestProblems(uint) ([]*model.ContestProblem, error)
	AddContestParticipants(uint, []uint) error
	AddContestAdmins(uint, []uint) error
	AddContestProblem(uint, uint) error
	DeleteContestProblem(uint, uint) error
	DeleteContestAdmin(uint, uint) error
	UpdateContest(*model.Contest) (*model.Contest, error)
	UpdateContestProblems(uint, []*model.ContestProblem) ([]*model.ContestProblem, error)
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

func (c *ContestMapperImpl) CreateContest(contest *model.Contest, problems []*model.ContestProblem) (*model.Contest, error) {
	tx := c.DB.Begin()
	//避免更新user
	user := contest.User
	contest.User = nil
	if err := tx.Create(contest).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, p := range problems {
		p.ContestId = contest.ID
		if err := tx.Create(p).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	/*未commit之前find不到
	contest, _ = c.FindContestById(contestId)
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	*/
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	contest.User = user
	return contest, nil
}

func (c *ContestMapperImpl) FindAllContests(pageNo int32, PageSize int32) ([]*model.Contest, int32, error) {
	limit, offset := util.CalLimitOffset(pageNo, PageSize)
	var count int32
	var contests []*model.Contest
	result := c.DB.Debug().
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

func (c *ContestMapperImpl) FindContestsByCondition(condition *SearchContestCondition, pageNo int32, pageSize int32) ([]*model.Contest, int32, error) {
	if condition == nil || (condition.Title == "" && condition.Status == 0 && condition.CreatorName == "") {
		return c.FindAllContests(pageNo, pageSize)
	}
	limit, offset := util.CalLimitOffset(pageNo, pageSize)
	result := c.DB.Model(&model.Contest{})
	var count int32
	var contests []*model.Contest
	now := time.Now()
	if condition.Status == contest_status.SCHEDULED {
		result = result.Where("start_time > ?", now)
	} else if condition.Status == contest_status.RUNNING {
		result = result.Where("start_time < ? and ? < end_time", now, now)
	} else if condition.Status == contest_status.ENDED {
		result = result.Where("end_time < ?", now)
	}
	if condition.Title != "" {
		result = result.Where("title like ?", fmt.Sprintf("%%%v%%", condition.Title))
	}
	if condition.CreatorName != "" {
		user, err := user_mapper.UserMapper.FindUserByUsername(condition.CreatorName)
		//没有此用户 直接返回
		if err != nil || user == nil {
			return contests, 0, nil
		}
		result = result.Where("user_id = ?", user.ID)
	}
	result = result.
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

func (c *ContestMapperImpl) AddContestProblem(contestId uint, problemId uint) error {
	contestProblem := &model.ContestProblem{
		ContestId: contestId,
		ProblemId: problemId,
	}
	result := c.DB.Create(contestProblem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *ContestMapperImpl) DeleteContestProblem(contestId uint, problemId uint) error {
	result := c.DB.
		Where("contest_id = ? and problem_id = ?", contestId, problemId).
		Delete(&model.ContestProblem{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *ContestMapperImpl) DeleteContestAdmin(contestId uint, userId uint) error {
	result := c.DB.
		Where("contest_id = ? and user_id = ?", contestId, userId).
		Delete(&model.ContestAdmin{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *ContestMapperImpl) UpdateContest(contest *model.Contest) (*model.Contest, error) {
	if contest.ID == 0 {
		return nil, fmt.Errorf("update contest need contest id")
	}
	user := contest.User
	contest.User = nil
	result := c.DB.Model(&model.Contest{}).Update(contest)
	if result.Error != nil {
		return nil, result.Error
	}
	contest.User = user
	return contest, nil
}

func (c *ContestMapperImpl) UpdateContestProblems(contestId uint, problems []*model.ContestProblem) ([]*model.ContestProblem, error) {
	tx := c.DB.Begin()
	if err := tx.Where("contest_id = ?", contestId).Delete(&model.ContestProblem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, p := range problems {
		p.ContestId = contestId
		if err := tx.Create(p).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return problems, nil
}

func (c *ContestMapperImpl) FindContestProblems(contestId uint) ([]*model.ContestProblem, error) {
	var contestProblems []*model.ContestProblem
	result := c.DB.
		Table("contest_problems").
		Where("contest_id = ?", contestId).
		Find(&contestProblems)
	if result.Error != nil {
		return nil, result.Error
	}
	return contestProblems, nil
}
