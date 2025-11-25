package handlers

import "github.com/AntonTsoy/review-pull-request-service/internal/service"

type PullRequestHandler struct {
	prService *service.PullRequestService
}

func newPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		prService: prService,
	}
}
