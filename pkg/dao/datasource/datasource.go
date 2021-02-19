package datasource

import (
	"fmt"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/contest_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/problem_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/submission_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/mapper/user_mapper"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type MysqlConf struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"db_name"`
	Local     string `yaml:"loc"`
	Charset   string `yaml:"charset"`
	ParseTime string `yaml:"parse_time"`
}

func loadConfig(path string) (conf *MysqlConf, err error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

var DB *gorm.DB

func ConnectDB(path string) error {
	conf, err := loadConfig(path)
	if err != nil {
		return err
	}
	DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&loc=%v&parseTime=%v", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName, conf.Charset, conf.Local, conf.ParseTime))
	if err != nil {
		return err
	}
	initMappers()
	migrateTables()
	return nil
}

func initMappers() {
	user_mapper.InitMapper(DB)
	submission_mapper.InitMapper(DB)
	problem_mapper.InitMapper(DB)
	contest_mapper.InitMapper(DB)
}

func migrateTables() {
	DB.AutoMigrate(
		&model.User{},
		&model.UserAuth{},
		&model.Role{},
		&model.Submission{},
		&model.SubmissionCode{},
		&model.CompileInfo{},
		&model.RawProblem{},
		&model.ProblemGroup{},
		&model.Problem{},
		&model.Contest{},
		&model.ContestProblem{},
		&model.ContestParticipant{},
		&model.ContestAdmin{},
	)
}
