package handlers

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
	"github.com/AntonTsoy/review-pull-request-service/internal/service"
)

type Handlers struct {
	api.ServerInterface

	*TeamHanler
	*UserHandler
	*PullRequestHandler
}

func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{
		TeamHanler: newTeamHandler(service.TeamService),
		UserHandler: newUserHandler(service.UserService),
		PullRequestHandler: newPullRequestHandler(service.PullRequestService),
	}
}
