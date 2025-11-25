package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Repository struct {
	TeamRepository *TeamRepository
	UserRepository *UserRepository
	ReviewRepository *ReviewRepository
	PullRequestRepository *PullRequestRepository
}

func NewRepository() *Repository {
	return &Repository{
		TeamRepository: newTeamRepository(),
		UserRepository: newUserRepository(),
		ReviewRepository: newReviewRepository(),
		PullRequestRepository: newPullRequestRepository(),
	}
}
