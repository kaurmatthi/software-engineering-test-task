package service

import (
	mock_repository "cruder/internal/mocks/repository"
	"cruder/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Given: A user is created with valid username, email, and full name
func TestCreateUser_Success(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "john_doe",
		Email:    "john@doe.ee",
		FullName: "John Doe",
	}

	mockRepo.EXPECT().Create(newUser).Return(newUser, nil).Times(1)

	// When: Calling create from user service
	createdUser, err := userService.Create(newUser)

	// Then: The result should be the created user and no error
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, newUser, createdUser, "expected created user to match input user")
}

// Given: A user is created with invalid username
func TestCreateUser_InvalidUsername_Fails(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "jo",
		Email:    "john@doe.ee",
		FullName: "John Doe",
	}

	mockRepo.EXPECT().Create(newUser).Return(newUser, nil).Times(0)

	// When: Calling create from user service
	createdUser, err := userService.Create(newUser)

	// Then: The result should be an ErrInvalidUsername error and nil user
	assert.ErrorIs(t, err, ErrInvalidUsername, "expected invalid username error")
	assert.Nil(t, createdUser, "expected no user to be returned")
}

// Given: A user is created with invalid email
func TestCreateUser_InvalidEmail_Fails(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "john_doe",
		Email:    "johndoe.ee",
		FullName: "John Doe",
	}

	mockRepo.EXPECT().Create(newUser).Return(newUser, nil).Times(0)

	// When: Calling create from user service
	createdUser, err := userService.Create(newUser)

	// Then: The result should be an ErrInvalidEmail error and nil user
	assert.ErrorIs(t, err, ErrInvalidEmail, "expected invalid email error")
	assert.Nil(t, createdUser, "expected no user to be returned")
}

// Given: A user is created with invalid full name
func TestCreateUser_InvalidFullName_Fails(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "john_doe",
		Email:    "john@doe.ee",
		FullName: "1John Doe",
	}

	mockRepo.EXPECT().Create(newUser).Return(newUser, nil).Times(0)

	// When: Calling create from user service
	createdUser, err := userService.Create(newUser)

	// Then: The result should be an ErrInvalidFullName error and nil user
	assert.ErrorIs(t, err, ErrInvalidFullName, "expected invalid email error")
	assert.Nil(t, createdUser, "expected no user to be returned")
}

// Given: A user is updated with valid username, email, and full name
func TestUpdateUser_Success(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "john_doe",
		Email:    "john@doe.ee",
		FullName: "John Doe",
	}

	mockRepo.EXPECT().Update(newUser).Return(newUser, nil).Times(1)

	// When: Calling create from user service
	createdUser, err := userService.Update(newUser)

	// Then: The result should be the created user and no error
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, newUser, createdUser, "expected created user to match input user")
}

// Given: A user is updated with invalid username
func TestUpdateUser_InvalidUsername_Fails(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "jo",
		Email:    "john@doe.ee",
		FullName: "John Doe",
	}

	mockRepo.EXPECT().Update(newUser).Return(newUser, nil).Times(0)

	// When: Calling create from user service
	createdUser, err := userService.Update(newUser)

	// Then: The result should be an ErrInvalidUsername error and nil user
	assert.ErrorIs(t, err, ErrInvalidUsername, "expected invalid username error")
	assert.Nil(t, createdUser, "expected no user to be returned")
}

// Given: A user is updated with invalid email
func TestUpdateUser_InvalidEmail_Fails(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "john_doe",
		Email:    "johndoe.ee",
		FullName: "John Doe",
	}

	mockRepo.EXPECT().Update(newUser).Return(newUser, nil).Times(0)

	// When: Calling create from user service
	createdUser, err := userService.Update(newUser)

	// Then: The result should be an ErrInvalidEmail error and nil user
	assert.ErrorIs(t, err, ErrInvalidEmail, "expected invalid email error")
	assert.Nil(t, createdUser, "expected no user to be returned")
}

// Given: A user is updated with invalid full name
func TestUpdateUser_InvalidFullName_Fails(t *testing.T) {
	// Setup: Create mock repository, service and user
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	newUser := &model.User{
		Username: "john_doe",
		Email:    "john@doe.ee",
		FullName: "1John Doe",
	}

	mockRepo.EXPECT().Update(newUser).Return(newUser, nil).Times(0)

	// When: Calling create from user service
	createdUser, err := userService.Update(newUser)

	// Then: The result should be an ErrInvalidFullName error and nil user
	assert.ErrorIs(t, err, ErrInvalidFullName, "expected invalid full name error")
	assert.Nil(t, createdUser, "expected no user to be returned")
}
