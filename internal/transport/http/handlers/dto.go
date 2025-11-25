package handlers

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/models"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
)

type ErrorMessage struct {
	Code    api.ErrorResponseErrorCode `json:"code"`
	Message string                     `json:"message"`
}

type TeamResponse struct {
	Team api.Team `json:"team"`
}

type UserGetReviewResponse struct {
	UserId       string                 `json:"user_id"`
	PullRequests []api.PullRequestShort `json:"pull_requests"`
}

type PullRequestResponse struct {
	Pr *api.PullRequest `json:"pr"`
}

type PullRequestAssignResponse struct {
	Pr         *api.PullRequest `json:"pr"`
	ReplacedBy string           `json:"replaced_by"`
}

func convertPRToAPI(pr *models.PullRequest) *api.PullRequest {
	return &api.PullRequest{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Title,
		AuthorId:          pr.AuthorID,
		Status:            api.PullRequestStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		MergedAt:          pr.MergedAt,
	}
}
