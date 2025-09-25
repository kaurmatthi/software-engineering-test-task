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

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.GetByUsername(username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := c.service.GetByID(id)
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	createdUser, err := c.service.Create(&user)
	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = c.service.Delete(id)

	if handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user model.User
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if id != user.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id in path and body do not match"})
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
		ctx.JSON(status, gin.H{"error": err.Error()})
		return true
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
