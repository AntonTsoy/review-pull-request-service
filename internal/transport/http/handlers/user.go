package handlers

import "github.com/AntonTsoy/review-pull-request-service/internal/service"

type UserHandler struct {
	userService *service.UserService
}

func newUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
