package controller

import (
	"net/http"
	"strconv"

	"cruder/internal/model"
	"cruder/internal/repository"
	"cruder/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Security ApiKeyAuth
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAll()
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserByUsername godoc
// @Summary Get user by username
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} model.User
// @Failure 404 {object} model.ErrorResponse "user not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Security ApiKeyAuth
// @Router /users/username/{username} [get]
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetByUsername(username)
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorResponse "invalid id"
// @Failure 404 {object} model.ErrorResponse "user not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Security ApiKeyAuth
// @Router /users/id/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	user, err := c.service.GetByID(id)
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.User true "User"
// @Success 201 {object} model.User
// @Failure 400 {object} model.ErrorResponse "invalid request body"
// @Failure 409 {object} model.ErrorResponse "user already exists"
// @Failure 409 {object} model.ErrorResponse "user with username/email already exists"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Security ApiKeyAuth
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}

	createdUser, err := c.service.Create(&user)
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} model.ErrorResponse "invalid id"
// @Failure 404 {object} model.ErrorResponse "user not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}
	err = c.service.Delete(id)

	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// UpdateUser godoc
// @Summary Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.User true "User"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorResponse "invalid id or body mismatch"
// @Failure 404 {object} model.ErrorResponse "user not found"
// @Failure 409 {object} model.ErrorResponse "user with username/email already exists"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user model.User
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}
	if id != user.ID {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "id in path and body do not match"})
		return
	}

	updatedUser, err := c.service.Update(&user)
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func handleError(ctx *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if status, ok := errToStatus[err]; ok {
		ctx.JSON(status, model.ErrorResponse{Error: err.Error()})
		return true
	}

	ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
	return true
}

var errToStatus = map[error]int{
	repository.ErrUserNotFound:          http.StatusNotFound,
	repository.ErrUserAlreadyExists:     http.StatusConflict,
	repository.ErrUsernameAlreadyExists: http.StatusConflict,
	repository.ErrEmailAlreadyExists:    http.StatusConflict,
	service.ErrInvalidEmail:             http.StatusBadRequest,
	service.ErrInvalidUsername:          http.StatusBadRequest,
	service.ErrInvalidFullName:          http.StatusBadRequest,
}
