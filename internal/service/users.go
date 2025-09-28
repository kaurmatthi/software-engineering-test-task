package service

import (
	"cruder/internal/model"
	"cruder/internal/repository"
	"errors"
	"regexp"
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
	if err := ValidateUser(*user); err != nil {
		return nil, err
	}
	return s.repo.Create(user)
}

func (s *userService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *userService) Update(user *model.User) (*model.User, error) {
	if err := ValidateUser(*user); err != nil {
		return nil, err
	}
	return s.repo.Update(user)
}

func ValidateUser(user model.User) error {
	if !emailRegex.MatchString(user.Email) {
		return ErrInvalidEmail
	}
	if !usernameRegex.MatchString(user.Username) {
		return ErrInvalidUsername
	}
	if !fullNameRegex.MatchString(user.FullName) {
		return ErrInvalidFullName
	}
	return nil
}

var emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
var usernameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{2,49}$`)
var fullNameRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z' -]{0,98}[A-Za-z]$`)

var ErrInvalidEmail = errors.New("invalid email format")
var ErrInvalidUsername = errors.New("invalid username format (3-50 chars, lowercase letters, numbers, underscores, starts with letter)")
var ErrInvalidFullName = errors.New("invalid full name format (2-100 chars, letters, spaces, apostrophes, hyphens, starts/ends with letter)")
