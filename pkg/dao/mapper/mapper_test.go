package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/language"
	"github.com/ecnuvj/vhoj_common/pkg/common/constants/remote_oj"
	"github.com/ecnuvj/vhoj_db/pkg/dao/datasource"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/problem_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/submission_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/user_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/jinzhu/gorm"
	"testing"
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
	}
	err := user_mapper.UserMapper.AddUser(user)
	if err != nil {
		fmt.Printf("add err: %v", err)
		return
	}
	fmt.Println(user.ID)
}

func TestUserMapperUpdateUser(t *testing.T) {
	connectDB()
	user := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
		UserAuth: &model.UserAuth{
			Password: "123456",
		},
		Email: "1486126243@qq.com",
	}
	//user_auth 会被新建
	err := user_mapper.UserMapper.UpdateUser(user)
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
	ret, err := submission_mapper.SubmissionMapper.AddSubmission(submission)
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
