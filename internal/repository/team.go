package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type TeamRepository struct {
	builder squirrel.StatementBuilderType
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *TeamRepository) Create(ctx context.Context, db DBTX, name string) (int, error) {
	sql, args, err := r.builder.
		Insert("teams").
		Columns("name").
		Values(name).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	var id int
	err = db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("execute query: %w", err)
	}

	return id, nil
}

func (r *TeamRepository) GetByName(ctx context.Context, db DBTX, name string) (int, error) {
	sql, args, err := r.builder.
		Select("id").
		From("teams").
		Where(squirrel.Eq{"name": name}).
		Limit(1).
		ToSql()
	
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	var id int
	err = db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperrors.ErrNotFound
		}
		return 0, fmt.Errorf("execute query: %w", err)
	}

	return id, nil
}
