package repository

import (
	"dryve/internal/datastruct"
	"dryve/internal/dto"

	"gorm.io/gorm"
)

type UserQuery interface {
	GetUser(id uint) (*datastruct.User, error)
	GetUserByEmail(email string) (*datastruct.User, error)
	CreateUser(user dto.RegisterRequest) (*datastruct.User, error)
	UpdateUser(user *datastruct.User) error
}

type userQuery struct {
	db *gorm.DB
}

func (d *dao) NewUserQuery() UserQuery {
	return &userQuery{d.db}
}

func (u *userQuery) GetUser(id uint) (*datastruct.User, error) {
	var user datastruct.User
	err := u.db.First(&user, id).Error
	return &user, err
}

func (u *userQuery) GetUserByEmail(email string) (*datastruct.User, error) {
	var user datastruct.User
	err := u.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (u *userQuery) CreateUser(user dto.RegisterRequest) (*datastruct.User, error) {
	var newUser datastruct.User
	newUser.FirstName = user.FirstName
	newUser.LastName = user.LastName
	newUser.Email = user.Email
	newUser.Password = user.Password
	err := u.db.Create(&newUser).Error
	return &newUser, err
}

func (u *userQuery) UpdateUser(user *datastruct.User) error {
	err := u.db.Save(user).Error
	return err
}
