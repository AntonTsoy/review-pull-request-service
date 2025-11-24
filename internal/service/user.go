package service

import (
	"context"

	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/models"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"
)

type UserService struct {
	db       *database.Database
	userRepo *repository.UserRepository
}

func newUserService(db *database.Database, userRepo *repository.UserRepository) *UserService {
	return &UserService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (*models.User, error) {
	user, err := s.userRepo.UpdateIsActive(ctx, s.db.Pool(), userID, active)
	if err != nil {
		return nil, err
	}

	// TODO: надо получить temaName

	return user, nil
}
