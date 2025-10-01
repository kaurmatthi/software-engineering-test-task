package controller

import (
	"bytes"
	mock_service "cruder/internal/mocks/service"
	"cruder/internal/model"
	"cruder/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupUserRouter(c *UserController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/users", c.GetAllUsers)
	r.GET("/users/username/:username", c.GetUserByUsername)
	r.GET("/users/id/:id", c.GetUserByID)
	r.POST("/users", c.CreateUser)
	r.DELETE("/users/:id", c.DeleteUser)
	r.PUT("/users/:id", c.UpdateUser)
	return r
}

func TestGetAllUsers_Success(t *testing.T) {
	// Given: service returns a list of users
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	users := []model.User{{ID: 1, Username: "john"}}
	mockSvc.EXPECT().GetAll().Return(users, nil)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: GET /users is called
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should contain the users
	assert.Equal(t, http.StatusOK, w.Code)
	var got []model.User
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, users, got)
}

func TestGetUserByUsername_Success(t *testing.T) {
	// Given: service returns a user for username "john_doe"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	expected := &model.User{ID: 1, Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	mockSvc.EXPECT().GetByUsername("john_doe").Return(expected, nil)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: GET /users/username/john_doe is called
	req, _ := http.NewRequest("GET", "/users/username/john_doe", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 200 with the expected user
	assert.Equal(t, http.StatusOK, w.Code)
	var got model.User
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, *expected, got)
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	// Given: service returns ErrUserNotFound
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	mockSvc.EXPECT().GetByUsername("missing").Return(nil, service.ErrUserNotFound)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: GET /users/username/missing is called
	req, _ := http.NewRequest("GET", "/users/username/missing", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 404 with error
	assert.Equal(t, http.StatusNotFound, w.Code)
	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "user not found", got.Error)
}

func TestGetUserByID_Success(t *testing.T) {
	// Given: service returns a user for ID 1
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	expected := &model.User{ID: 1, Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	mockSvc.EXPECT().GetByID(int64(1)).Return(expected, nil)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: GET /users/id/1 is called
	req, _ := http.NewRequest("GET", "/users/id/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 200 with the expected user
	assert.Equal(t, http.StatusOK, w.Code)
	var got model.User
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, *expected, got)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	// Given: invalid ID in path
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: GET /users/id/abc is called
	req, _ := http.NewRequest("GET", "/users/id/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 400 with error
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "invalid id", got.Error)
}

func TestCreateUser_Success(t *testing.T) {
	// Given: valid user and service returns created user
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	input := &model.User{Username: "john", Email: "john@doe.ee", FullName: "John Doe"}
	created := &model.User{ID: 1, Username: "john", Email: "john@doe.ee", FullName: "John Doe"}
	mockSvc.EXPECT().Create(input).Return(created, nil)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	body, _ := json.Marshal(input)

	// When: POST /users is called
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 201 with created user
	assert.Equal(t, http.StatusCreated, w.Code)
	var got model.User
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, *created, got)
}

func TestCreateUser_InvalidBody(t *testing.T) {
	// Given: malformed JSON body
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: POST /users with invalid JSON
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte("{bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 400 with error
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "invalid request body", got.Error)
}

func TestDeleteUser_NotFound(t *testing.T) {
	// Given: service returns ErrUserNotFound
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	mockSvc.EXPECT().Delete(int64(99)).Return(service.ErrUserNotFound)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	// When: DELETE /users/99 is called
	req, _ := http.NewRequest("DELETE", "/users/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 404 with error
	assert.Equal(t, http.StatusNotFound, w.Code)
	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "user not found", got.Error)
}

func TestUpdateUser_IDMismatch(t *testing.T) {
	// Given: ID in path and body do not match
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	user := model.User{ID: 2, Username: "john", Email: "john@doe.ee", FullName: "John Doe"}
	body, _ := json.Marshal(user)

	// When: PUT /users/1 with body ID 2
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 400 with mismatch error
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "id in path and body do not match", got.Error)
}

func TestUpdateUser_Success(t *testing.T) {
	// Given: valid update request and service returns updated user
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mock_service.NewMockUserService(ctrl)
	input := &model.User{ID: 1, Username: "john", Email: "john@doe.ee", FullName: "John Doe"}
	updated := &model.User{ID: 1, Username: "johnny", Email: "johnny@doe.ee", FullName: "Johnny Doe"}
	mockSvc.EXPECT().Update(input).Return(updated, nil)

	controller := NewUserController(mockSvc)
	router := setupUserRouter(controller)

	body, _ := json.Marshal(input)

	// When: PUT /users/1 is called with valid body
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: response should be 200 with updated user
	assert.Equal(t, http.StatusOK, w.Code)
	var got model.User
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, *updated, got)
}

func TestHandleError_UserNotFound(t *testing.T) {
	// Given: a context and ErrUserNotFound
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrUserNotFound
	result := handleError(c, service.ErrUserNotFound)

	// Then: should return true and response with 404 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "user not found", got.Error)
}

func TestHandleError_UserAlreadyExists(t *testing.T) {
	// Given: a context and ErrUserAlreadyExists
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrUserAlreadyExists
	result := handleError(c, service.ErrUserAlreadyExists)

	// Then: should return true and response with 409 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusConflict, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "user already exists", got.Error)
}

func TestHandleError_UsernameAlreadyExists(t *testing.T) {
	// Given: a context and ErrUsernameAlreadyExists
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrUsernameAlreadyExists
	result := handleError(c, service.ErrUsernameAlreadyExists)

	// Then: should return true and response with 409 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusConflict, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "username already exists", got.Error)
}

func TestHandleError_EmailAlreadyExists(t *testing.T) {
	// Given: a context and ErrEmailAlreadyExists
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrEmailAlreadyExists
	result := handleError(c, service.ErrEmailAlreadyExists)

	// Then: should return true and response with 409 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusConflict, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "email already exists", got.Error)
}

func TestHandleError_InvalidEmail(t *testing.T) {
	// Given: a context and ErrInvalidEmail
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrInvalidEmail
	result := handleError(c, service.ErrInvalidEmail)

	// Then: should return true and response with 400 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Contains(t, got.Error, "invalid email")
}

func TestHandleError_InvalidUsername(t *testing.T) {
	// Given: a context and ErrInvalidUsername
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrInvalidUsername
	result := handleError(c, service.ErrInvalidUsername)

	// Then: should return true and response with 400 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Contains(t, got.Error, "invalid username")
}

func TestHandleError_InvalidFullName(t *testing.T) {
	// Given: a context and ErrInvalidFullName
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// When: calling handleError with ErrInvalidFullName
	result := handleError(c, service.ErrInvalidFullName)

	// Then: should return true and response with 400 + correct error message
	assert.True(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Contains(t, got.Error, "invalid full name")
}

func TestHandleError_UnknownError(t *testing.T) {
	// Given: a context and an unknown error
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	unknown := assert.AnError

	// When: calling handleError with unknown error
	result := handleError(c, unknown)

	// Then: should return true and response with 500 + generic message
	assert.True(t, result)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var got model.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "internal server error", got.Error)
}
