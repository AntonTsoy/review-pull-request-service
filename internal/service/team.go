package service

import (
	"context"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/api"
	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"
	"github.com/jackc/pgx/v5"
)

type TeamRepository interface {
	Create(ctx context.Context, db repository.DBTX, name string) (int, error)
	GetByName(ctx context.Context, db repository.DBTX, name string) (int, error)
}

type TeamService struct {
	db       *database.Database
	teamRepo TeamRepository
}

func NewTeamService(db *database.Database, teamRepo TeamRepository) *TeamService {
	return &TeamService{
		db:       db,
		teamRepo: teamRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team api.Team) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	id, err := s.teamRepo.GetByName(ctx, tx, team.TeamName)
	if err == nil {
		return apperrors.ErrTeamExists
	}

	id, err = s.teamRepo.Create(ctx, tx, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	// TODO: create users/members

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
