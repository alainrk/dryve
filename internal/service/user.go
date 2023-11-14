package service

import (
	"dryve/internal/datastruct"
	"dryve/internal/dto"
	"dryve/internal/repository"
	"dryve/internal/utils"
)

type UserService interface {
	GetUser(id uint) (*datastruct.User, error)
	GetUserByEmail(email string) (*datastruct.User, error)
	CreateUser(dto.RegisterRequest) (*datastruct.User, error)
	SetEmailConfirmationCode(userId uint) (string, error)
	VerifyUser(userId uint) error
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (s *userService) GetUser(id uint) (*datastruct.User, error) {
	return s.dao.NewUserQuery().GetUser(id)
}

func (s *userService) GetUserByEmail(email string) (*datastruct.User, error) {
	return s.dao.NewUserQuery().GetUserByEmail(email)
}

func (s *userService) CreateUser(user dto.RegisterRequest) (*datastruct.User, error) {
	return s.dao.NewUserQuery().CreateUser(user)
}

func (s *userService) SetEmailConfirmationCode(userId uint) (string, error) {
	code := utils.RandSeq(6)
	user, err := s.dao.NewUserQuery().GetUser(userId)
	if err != nil {
		return code, err
	}
	user.EmailCode = code
	err = s.dao.NewUserQuery().UpdateUser(user)
	return code, err
}

func (s *userService) VerifyUser(userId uint) error {
	user, err := s.dao.NewUserQuery().GetUser(userId)
	if err != nil {
		return err
	}
	user.Verified = true
	err = s.dao.NewUserQuery().UpdateUser(user)
	return err
}
