package service

import (
	"context"

	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
)

type UserService struct {
	db       *database.Database
	userRepo *repository.UserRepository
	teamRepo *repository.TeamRepository
}

func newUserService(db *database.Database, userRepo *repository.UserRepository, teamRepo *repository.TeamRepository) *UserService {
	return &UserService{
		db:       db,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (*api.User, error) {
	user, err := s.userRepo.UpdateIsActive(ctx, s.db.Pool(), userID, active)
	if err != nil {
		return nil, err
	}

	teamName, err := s.teamRepo.GetByID(ctx, s.db.Pool(), user.TeamID)
	if err != nil {
		return nil, err
	}

	return &api.User{
		UserId:   user.ID,
		Username: user.Name,
		IsActive: user.IsActive,
		TeamName: teamName,
	}, nil
}
