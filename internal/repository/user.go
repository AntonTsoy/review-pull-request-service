package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/api"
	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/Masterminds/squirrel"
)

type UserRepository struct {
	builder squirrel.StatementBuilderType
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) Create(ctx context.Context, db DBTX, teamID int, user *api.TeamMember) error {
	sql, args, err := r.builder.
		Insert("users").
		Columns("id", "name", "is_active", "team_id").
		Values(user.UserId, user.Username, user.IsActive, teamID).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	if _, err = db.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByTeamID(ctx context.Context, db DBTX, teamID int) ([]api.TeamMember, error) {
	sql, args, err := r.builder.
		Select("id", "name", "is_active").
		From("users").
		Where(squirrel.Eq{"team_id": teamID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var teamMembers []api.TeamMember
	for rows.Next() {
		var member api.TeamMember
		err := rows.Scan(&member.UserId, &member.Username, &member.IsActive)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		teamMembers = append(teamMembers, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teamMembers, nil
}

func (r *UserRepository) Update(ctx context.Context, db DBTX, id string, name string, teamID int) error {
	sql, args, err := r.builder.
		Update("users").
		Set("name", name).
		Set("team_id", teamID).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	cmdTag, err := db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *UserRepository) SetActive(ctx context.Context, db DBTX, id string, active bool) (models.User, error) {
	sql, args, err := r.builder.
		Update("users").
		Set("is_active", active).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, name, is_active, team_id").
		ToSql()

	if err != nil {
		return models.User{}, fmt.Errorf("build query: %w", err)
	}

	var user models.User
	err = db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.IsActive,
		&user.TeamID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, apperrors.ErrNotFound
		}
		return models.User{}, fmt.Errorf("execute query and scan result: %w", err)
	}

	return user, nil
}
