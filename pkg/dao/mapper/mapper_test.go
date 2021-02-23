package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/language"
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/remote_oj"
	"github.com/ecnuvj/vhoj_db/pkg/dao/datasource"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/contest_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/problem_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/submission_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/user_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

func connectDB() {
	err := datasource.ConnectDB("../../../config/mysql.yaml")
	if err != nil {
		fmt.Printf("connect err: %v", err)
		return
	}
}

func TestUserMapperAddUser(t *testing.T) {
	connectDB()
	user := &model.User{
		UserAuth: &model.UserAuth{
			Password: "123456",
		},
		Nickname: "bqx",
	}
	retUser, err := user_mapper.UserMapper.AddUser(user)
	if err != nil {
		fmt.Printf("add err: %v", err)
		return
	}
	str, _ := json.Marshal(retUser)
	fmt.Println(string(str))
}

func TestUserMapperUpdateUser(t *testing.T) {
	connectDB()
	user := &model.User{
		Model: gorm.Model{
			ID: 9,
		},
		UserAuth: &model.UserAuth{
			Model: gorm.Model{
				ID: 1,
			},
			Password: "123456456",
		},
		Email: "1486126243@qq.com",
	}
	//user_auth 会被新建 (因为未指定user_auth的主键)
	_, err := user_mapper.UserMapper.UpdateUser(user)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestUserMapperRelatedGetUser(t *testing.T) {
	connectDB()
	user := &model.User{
		Model: gorm.Model{ID: 1},
	}
	//var userAuth model.UserAuth
	datasource.DB.Model(user).Find(user).Related(&user.UserAuth)
	str, _ := json.Marshal(user)
	fmt.Println(string(str))
}

func TestSubmissionMapperAddSubmission(t *testing.T) {
	connectDB()
	submission := &model.Submission{
		SubmissionCode: &model.SubmissionCode{
			SourceCode: "I2luY2x1ZGUgPHN0ZGlvLmg%2BCgppbnQgbWFpbigpCnsKICAgIGludCBhLGI7CiAgICBzY2FuZigiJWQgJWQiLCZhLCAmYik7CiAgICBwcmludGYoIiVkXG4iLGErYik7CiAgICByZXR1cm4gMDsKfQ%3D%3D",
			CodeLength: 0,
		},
		ProblemId:  0,
		UserId:     0,
		Result:     0,
		TimeCost:   0,
		MemoryCost: 0,
		Language:   language.CPP,
		ContestId:  0,
		RealRunId:  "",
	}
	ret, err := submission_mapper.SubmissionMapper.AddOrModifySubmission(submission)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(ret)
	fmt.Println(string(str))
}

func TestSubmissionMapperFindSubmission(t *testing.T) {
	connectDB()
	submission, err := submission_mapper.SubmissionMapper.FindSubmissionById(3)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(submission)
	fmt.Println(string(str))
}

func TestProblemMapperAddOrModifyRawProblem(t *testing.T) {
	connectDB()
	rawProblem := &model.RawProblem{
		Title:           "problem 1000",
		RemoteOJ:        remote_oj.HDU,
		RemoteProblemId: "1000",
	}
	problem_mapper.ProblemMapper.AddOrModifyRawProblem(rawProblem)
}

func TestProblemMapperImpl_FindGroupProblemsById(t *testing.T) {
	connectDB()
	result, err := problem_mapper.ProblemMapper.FindGroupProblemsById(1)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(result)
	fmt.Println(string(str))
}

func TestProblemMapperImpl_FindAllProblems(t *testing.T) {
	connectDB()
	result, count, err := problem_mapper.ProblemMapper.FindAllProblems(1, 1)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(result)
	fmt.Printf("count: %v,result: %v", count, string(str))
}

func TestSubmissionMapperImpl_UpdateSubmissionById(t *testing.T) {
	connectDB()
	submission := &model.Submission{
		Model:      gorm.Model{ID: 5},
		TimeCost:   5,
		MemoryCost: 5,
		RealRunId:  "11111",
	}
	submission, err := submission_mapper.SubmissionMapper.UpdateSubmissionById(submission)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(submission)
	fmt.Println(string(str))
}

func TestSubmissionMapperImpl_ResetSubmissionById(t *testing.T) {
	connectDB()
	err := submission_mapper.SubmissionMapper.ResetSubmissionById(5)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestSubmissionMapperImpl_AddOrModifySubmissionById(t *testing.T) {
	connectDB()
	submission, err := submission_mapper.SubmissionMapper.FindSubmissionById(7)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	str, _ := json.Marshal(submission)
	fmt.Println(string(str))
	submission, err = submission_mapper.SubmissionMapper.AddOrModifySubmission(submission)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	str, _ = json.Marshal(submission)
	fmt.Println(string(str))
}

func TestProblemMapperImpl_AddProblemSubmittedCountById(t *testing.T) {
	connectDB()
	err := problem_mapper.ProblemMapper.AddProblemSubmittedCountById(1)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}

func TestProblemMapperImpl_AddProblemAcceptedCountById(t *testing.T) {
	connectDB()
	err := problem_mapper.ProblemMapper.AddProblemAcceptedCountById(1)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}

func TestProblemMapperImpl_GetProblemById(t *testing.T) {
	connectDB()
	problem, err := problem_mapper.ProblemMapper.FindProblemById(2)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(problem)
	fmt.Println(string(str))
}

func TestProblemMapperImpl_SearchProblemByCondition(t *testing.T) {
	connectDB()
	problems, _, err := problem_mapper.ProblemMapper.SearchProblemByCondition(&problem_mapper.ProblemSearchParam{
		Title:     "pen",
		ProblemId: 0,
	}, 1, 1)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(problems)
	fmt.Printf("len: %v\n", len(problems))
	fmt.Println(string(str))
}

func TestContestMapperImpl_CreateContest(t *testing.T) {
	connectDB()
	contest := &model.Contest{
		Title:       "first contest",
		Description: "hello world",
		UserId:      1,
		ProblemNum:  2,
		ProblemIds:  []uint{2, 3},
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour * 5),
	}
	contest, err := contest_mapper.ContestMapper.CreateContest(contest)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(contest)
	fmt.Println(string(str))
}

func TestContestMapperImpl_FindAllContests(t *testing.T) {
	connectDB()
	contests, _, err := contest_mapper.ContestMapper.FindAllContests(1, 5)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(contests)
	fmt.Println(string(str))
}

func TestContestMapperImpl_FindContestById(t *testing.T) {
	connectDB()
	contest, err := contest_mapper.ContestMapper.FindContestById(4)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(contest)
	fmt.Println(string(str))
}

func TestContestMapperImpl_AddContestParticipants(t *testing.T) {
	connectDB()
	err := contest_mapper.ContestMapper.AddContestParticipants(4, []uint{1, 2, 3, 4})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestContestMapperImpl_AddContestAdmins(t *testing.T) {
	connectDB()
	err := contest_mapper.ContestMapper.AddContestAdmins(4, []uint{1, 2, 3, 4})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestUserMapperImpl_FindUsersByIds(t *testing.T) {
	connectDB()
	users, err := user_mapper.UserMapper.FindUsersByIds([]uint{5, 6, 7})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(users)
	fmt.Println(string(str))
}

func TestUserMapperImpl_FindUserByUsername(t *testing.T) {
	connectDB()
	user, err := user_mapper.UserMapper.FindUserByUsername("")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(user)
	fmt.Println(string(str))
}

func TestUserMapperImpl_AddUserRoleByRoleName(t *testing.T) {
	connectDB()
	err := user_mapper.UserMapper.AddUserRoleByRoleName(2, "normal")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestUserMapperImpl_UpdateUserRoles(t *testing.T) {
	connectDB()
	err := user_mapper.UserMapper.UpdateUserRoles(12, []*model.Role{
		{Model: gorm.Model{ID: 7}},
		{Model: gorm.Model{ID: 8}},
		{Model: gorm.Model{ID: 9}},
	})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestUserMapperImpl_DeleteUserById(t *testing.T) {
	connectDB()
	err := user_mapper.UserMapper.DeleteUserById(10)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
}

func TestUserMapperImpl_FindAllUsers(t *testing.T) {
	connectDB()
	users, count, err := user_mapper.UserMapper.FindAllUsers(1, 5)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Printf("count: %v\n", count)
	str, _ := json.Marshal(users)
	fmt.Println(string(str))
}
