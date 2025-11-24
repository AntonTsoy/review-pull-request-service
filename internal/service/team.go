package service

import (
	"context"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"

	"github.com/jackc/pgx/v5"
)

type TeamService struct {
	db       *database.Database
	teamRepo *repository.TeamRepository
	userRepo *repository.UserRepository
}

func newTeamService(db *database.Database, teamRepo *repository.TeamRepository, userRepo *repository.UserRepository) *TeamService {
	return &TeamService{
		db:       db,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team api.Team) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	teamID, err := s.teamRepo.GetByName(ctx, tx, team.TeamName)
	if err == nil {
		return apperrors.ErrTeamExists
	}
	_ = teamID

	teamID, err = s.teamRepo.Create(ctx, tx, team.TeamName)
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		var flag bool
		flag, err = s.userRepo.Exists(ctx, tx, member.UserId)
		if err == nil && flag {
			if err = s.userRepo.Update(ctx, tx, teamID, &member); err != nil {
				return err
			}
			continue
		}

		if err = s.userRepo.Create(ctx, tx, teamID, &member); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) ([]api.TeamMember, error) {
	teamID, err := s.teamRepo.GetByName(ctx, s.db.Pool(), teamName)
	if err != nil {
		return nil, err
	}

	members, err := s.userRepo.GetByTeamID(ctx, s.db.Pool(), teamID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
