package problem_mapper

import (
	"encoding/json"
	"fmt"
	"github.com/ecnuvj/vhoj_db/pkg/dao/datasource"
	"testing"
)

func connectDB() {
	err := datasource.ConnectDB("../../../config/mysql.yaml")
	if err != nil {
		fmt.Printf("connect err: %v", err)
		return
	}
}

func TestProblemMapperImpl_FindGroupProblemsById(t *testing.T) {
	connectDB()
	result, err := ProblemMapper.FindGroupProblemsById(1)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	str, _ := json.Marshal(result)
	fmt.Println(string(str))
}
