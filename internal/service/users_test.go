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

	// When: calling create from service layer
	createdUser, err := userService.Create(newUser)

	// Then: Expected result
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

	// When: calling create from service layer
	createdUser, err := userService.Create(newUser)

	// Then: Expected result
	if err == nil {
		t.Fatalf("expected error, got no error")
	}
	if createdUser == newUser {
		t.Fatalf("expected no user, got %v", createdUser)
	}
}
