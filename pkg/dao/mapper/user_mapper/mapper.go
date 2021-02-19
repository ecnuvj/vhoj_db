package user_mapper

import (
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/jinzhu/gorm"
)

type IUserMapper interface {
	AddUser(*model.User) error
	UpdateUser(*model.User) error
	FindUsersByIds([]uint) ([]*model.User, error)
}

var UserMapper IUserMapper

func InitMapper(db *gorm.DB) {
	UserMapper = &UserMapperImpl{DB: db}
}

type UserMapperImpl struct {
	DB *gorm.DB
}

func (u *UserMapperImpl) AddUser(user *model.User) error {
	result := u.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserMapperImpl) UpdateUser(user *model.User) error {
	result := u.DB.Model(&model.User{}).Update(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserMapperImpl) FindUsersByIds(userIds []uint) ([]*model.User, error) {
	var users []*model.User
	result := u.DB.
		Model(&model.User{}).
		Find(&users, userIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
