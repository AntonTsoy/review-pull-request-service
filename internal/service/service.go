package service

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"
)

type Service struct {
	TeamService *TeamService
	UserService *UserService
	PullRequestService *PullRequestService
}

func NewService(db *database.Database, repo *repository.Repository) *Service {
	return &Service{
		TeamService: newTeamService(db, repo.TeamRepository, repo.UserRepository),
		UserService: newUserService(db, repo.UserRepository, repo.TeamRepository),
		PullRequestService: newPullRequestService(db, repo.PullRequestRepository, repo.UserRepository, repo.ReviewRepository),
	}
}
