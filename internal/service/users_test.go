package service

import (
	mock_repository "cruder/internal/mocks/repository"
	"cruder/internal/model"
	"testing"

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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdUser != newUser {
		t.Fatalf("expected user %v, got %v", newUser, createdUser)
	}
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
	if err != ErrInvalidUsername {
		t.Fatalf("expected invalid username error, got wrong or no error")
	}
	if createdUser != nil {
		t.Fatalf("expected no user, got %v", createdUser)
	}
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
	if err != ErrInvalidEmail {
		t.Fatalf("expected invalid email error, got wrong or no error")
	}
	if createdUser != nil {
		t.Fatalf("expected no user, got %v", createdUser)
	}
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
	if err != ErrInvalidFullName {
		t.Fatalf("expected invalid full name error, got wrong or no error")
	}
	if createdUser != nil {
		t.Fatalf("expected no user, got %v", createdUser)
	}
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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdUser != newUser {
		t.Fatalf("expected user %v, got %v", newUser, createdUser)
	}
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
	if err != ErrInvalidUsername {
		t.Fatalf("expected invalid username error, got wrong or no error")
	}
	if createdUser != nil {
		t.Fatalf("expected no user, got %v", createdUser)
	}
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
	if err != ErrInvalidEmail {
		t.Fatalf("expected invalid email error, got wrong or no error")
	}
	if createdUser != nil {
		t.Fatalf("expected no user, got %v", createdUser)
	}
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
	if err != ErrInvalidFullName {
		t.Fatalf("expected invalid full name error, got wrong or no error")
	}
	if createdUser != nil {
		t.Fatalf("expected no user, got %v", createdUser)
	}
}
