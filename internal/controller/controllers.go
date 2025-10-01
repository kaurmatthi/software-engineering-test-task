package controller

import "cruder/internal/service"

type Controller struct {
	Users  *UserController
	Health *HealthController
}

func NewController(services *service.Service) *Controller {
	return &Controller{
		Users:  NewUserController(services.Users),
		Health: NewHealthController(),
	}
}
