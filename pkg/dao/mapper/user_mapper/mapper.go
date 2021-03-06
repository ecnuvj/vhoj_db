package user_mapper

import (
	"fmt"
	"github.com/ecnuvj/vhoj_db/pkg/dao/model"
	"github.com/ecnuvj/vhoj_db/pkg/util"
	"github.com/jinzhu/gorm"
)

type IUserMapper interface {
	AddUser(*model.User) (*model.User, error)
	AddUserRoleByRoleName(uint, string) (*model.Role, error)
	AddUserRoleByRoleId(uint, uint) error
	AddUserSubmitCountById(uint) error
	AddUserAcceptCountById(uint) error
	UpdateUser(*model.User) (*model.User, error)
	UpdateUserRoles(uint, []*model.Role) (*model.User, error)
	FindUsersByIds([]uint) ([]*model.User, error)
	FindUserById(uint) (*model.User, error)
	FindUserByUsername(string) (*model.User, error)
	FindUserRolesById(uint) ([]*model.Role, error)
	FindAllUsers(int32, int32) ([]*model.User, int32, error)
	FindRoleList() ([]*model.Role, error)
	DeleteUserById(uint) error
}

var UserMapper IUserMapper

func InitMapper(db *gorm.DB) {
	UserMapper = &UserMapperImpl{DB: db}
}

type UserMapperImpl struct {
	DB *gorm.DB
}

func (u *UserMapperImpl) AddUser(user *model.User) (*model.User, error) {
	tx := u.DB.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for i, r := range user.Roles {
		if rol, err := u.AddUserRoleByRoleName(user.ID, r.RoleName); err != nil {
			tx.Rollback()
			return nil, err
		} else {
			user.Roles[i] = rol
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return user, nil
}

//只更新用户信息 不更新角色
func (u *UserMapperImpl) UpdateUser(user *model.User) (*model.User, error) {
	if user.UserAuth != nil && user.UserAuth.ID == 0 {
		user.UserAuth = nil
	}
	result := u.DB.Model(&model.User{}).Update(user)
	if result.Error != nil {
		return nil, result.Error
	}
	retUser, err := u.FindUserById(user.ID)
	if err != nil {
		return nil, err
	}
	return retUser, nil
}

//contest admins participants
func (u *UserMapperImpl) FindUsersByIds(userIds []uint) ([]*model.User, error) {
	var users []*model.User
	result := u.DB.
		Model(&model.User{}).
		Find(&users, userIds)
	if result.Error != nil {
		return nil, result.Error
	}
	for i, user := range users {
		users[i].Roles, _ = u.FindUserRolesById(user.ID)
	}
	return users, nil
}

func (u *UserMapperImpl) FindUserByUsername(username string) (*model.User, error) {
	user := &model.User{UserAuth: &model.UserAuth{}}
	result := u.DB.
		Model(user).
		Where("nickname = ?", username).
		Find(user)
	if result.Error != nil {
		return nil, result.Error
	}
	result = u.DB.Model(user).Related(user.UserAuth)
	if result.Error != nil {
		return nil, result.Error
	}
	roles, _ := u.FindUserRolesById(user.ID)
	user.Roles = roles
	return user, nil
}

func (u *UserMapperImpl) FindUserById(userId uint) (*model.User, error) {
	user := &model.User{
		Model:    gorm.Model{ID: userId},
		UserAuth: &model.UserAuth{},
	}
	result := u.DB.Model(user).First(user).Related(user.UserAuth)
	if result.Error != nil {
		return nil, result.Error
	}
	roles, _ := u.FindUserRolesById(userId)
	user.Roles = roles
	return user, nil
}

func (u *UserMapperImpl) AddUserRoleByRoleName(userId uint, roleName string) (*model.Role, error) {
	role := &model.Role{
		RoleName: roleName,
	}
	if err := u.DB.Model(&model.Role{}).Where("role_name = ?", roleName).First(role).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			u.DB.Create(role)
		}
	}
	err := u.AddUserRoleByRoleId(userId, role.ID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (u *UserMapperImpl) AddUserRoleByRoleId(userId uint, roleId uint) error {
	userRole := &model.UserRole{
		UserId: userId,
		RoleId: roleId,
	}
	result := u.DB.Create(userRole)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserMapperImpl) FindUserRolesById(userId uint) ([]*model.Role, error) {
	var userRoles []*model.UserRole
	result := u.DB.Model(&model.UserRole{}).Where("user_id = ?", userId).Find(&userRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	roleIds := make([]uint, len(userRoles))
	for i, r := range userRoles {
		roleIds[i] = r.RoleId
	}
	var roles []*model.Role
	result = u.DB.Model(&model.Role{}).Find(&roles, roleIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

//role id 必须给到
func (u *UserMapperImpl) UpdateUserRoles(userId uint, roles []*model.Role) (*model.User, error) {
	tx := u.DB.Begin()
	if err := tx.Where("user_id = ?", userId).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, r := range roles {
		userRole := &model.UserRole{
			UserId: userId,
			RoleId: r.ID,
		}
		if err := tx.Create(userRole).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}
	return u.FindUserById(userId)
}

func (u *UserMapperImpl) DeleteUserById(userId uint) error {
	tx := u.DB.Begin()
	user := &model.User{
		Model:    gorm.Model{ID: userId},
		UserAuth: &model.UserAuth{},
	}
	if err := tx.First(user).Related(user.UserAuth).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(user.UserAuth).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("user_id = ?", userId).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (u *UserMapperImpl) FindAllUsers(pageNo int32, pageSize int32) ([]*model.User, int32, error) {
	limit, offset := util.CalLimitOffset(pageNo, pageSize)
	var count int32
	var users []*model.User
	result := u.DB.
		Model(&model.User{}).
		Count(&count).
		Preload("UserAuth").
		Limit(limit).
		Offset(offset).
		Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	for i, user := range users {
		users[i].Roles, _ = u.FindUserRolesById(user.ID)
	}
	return users, count, nil
}

func (u *UserMapperImpl) AddUserSubmitCountById(userId uint) error {
	if userId <= 0 {
		return fmt.Errorf("user id is incorrect")
	}
	result := u.DB.
		Model(&model.User{Model: gorm.Model{ID: userId}}).
		Update("submitted", gorm.Expr("submitted + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserMapperImpl) AddUserAcceptCountById(userId uint) error {
	if userId <= 0 {
		return fmt.Errorf("problem id is incorrect")
	}
	result := u.DB.
		Model(&model.User{Model: gorm.Model{ID: userId}}).
		Update("accepted", gorm.Expr("accepted + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserMapperImpl) FindRoleList() ([]*model.Role, error) {
	var roles []*model.Role
	result := u.DB.
		Model(&model.Role{}).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
