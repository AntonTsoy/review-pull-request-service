package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/Masterminds/squirrel"
)

type UserRepository struct {
	builder squirrel.StatementBuilderType
}

func newUserRepository() *UserRepository {
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

func (r *UserRepository) Exists(ctx context.Context, db DBTX, id string) (bool, error) {
	sql, args, err := r.builder.
		Select("1").
		From("users").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("build query: %w", err)
	}

    var exists bool
    err = db.QueryRow(ctx, sql, args...).Scan(&exists)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return false, nil
        }
        return false, fmt.Errorf("query row: %w", err)
    }

    return exists, nil
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

func (r *UserRepository) GetActiveTeammates(ctx context.Context, db DBTX, exceptID string) ([]api.TeamMember, error) {
	sql, args, err := r.builder.
		Select("u1.id", "u1.name", "u1.is_active").
		From("users u1").
		Join("users u2 ON u1.team_id = u2.team_id").
		Where(squirrel.Eq{"u2.user_id": exceptID}).
		Where(squirrel.NotEq{"u1.user_id": exceptID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var teammates []api.TeamMember
	for rows.Next() {
		var member api.TeamMember
		err := rows.Scan(&member.UserId, &member.Username, &member.IsActive)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		teammates = append(teammates, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teammates, nil
}

func (r *UserRepository) Update(ctx context.Context, db DBTX, teamID int, user *api.TeamMember) error {
	sql, args, err := r.builder.
		Update("users").
		Set("name", user.Username).
		Set("team_id", teamID).
		Set("is_active", user.IsActive).
		Where(squirrel.Eq{"id": user.UserId}).
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

func (r *UserRepository) UpdateIsActive(ctx context.Context, db DBTX, id string, active bool) (*models.User, error) {
	sql, args, err := r.builder.
		Update("users").
		Set("is_active", active).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, name, is_active, team_id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	user := &models.User{}
	err = db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.IsActive,
		&user.TeamID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("execute query and scan result: %w", err)
	}

	return user, nil
}
