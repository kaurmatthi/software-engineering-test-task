package service

import (
	"cruder/internal/model"
	"cruder/internal/repository"
)

type UserService interface {
	GetAll() ([]model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Delete(id int64) error
	Update(user *model.User) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetByUsername(username string) (*model.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *userService) GetByID(id int64) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) Create(user *model.User) (*model.User, error) {
	// Validate user
	return s.repo.Create(user)
}

func (s *userService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *userService) Update(user *model.User) (*model.User, error) {
	// Validate user
	return s.repo.Update(user)
}
