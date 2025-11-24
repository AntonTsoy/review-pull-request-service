package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type PullRequestRepository struct {
	builder squirrel.StatementBuilderType
}

func newPullRequestRepository() *PullRequestRepository {
	return &PullRequestRepository{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PullRequestRepository) Create(ctx context.Context, db DBTX, pr *models.PullRequest) error {
	sql, args, err := r.builder.
		Insert("pull_requests").
		Columns("id", "title", "author_id").
		Values(pr.ID, pr.Title, pr.AuthorID).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create pr query: %w", err)
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("exec create pr: %w", err)
	}

	return nil
}

func (r *PullRequestRepository) GetByID(ctx context.Context, db DBTX, id string) (*models.PullRequest, error) {
	sql, args, err := r.builder.
		Select("id", "title", "author_id", "status", "merged_at").
		From("pull_requests").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var pr models.PullRequest
	err = db.QueryRow(ctx, sql, args...).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.MergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("scan pr: %w", err)
	}

	return &pr, nil
}

func (r *PullRequestRepository) Exists(ctx context.Context, db DBTX, id string) (bool, error) {
	sql, args, err := r.builder.
		Select("1").
		From("pull_requests").
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

func (r *PullRequestRepository) UpdateMergeStatus(ctx context.Context, db DBTX, id string) error {
	sql, args, err := r.builder.
		Update("pull_requests").
		Set("status", models.StatusMerged).
		Set("merged_at", "NOW()").
		Where(squirrel.Eq{"id": id, "status": models.StatusOpen}).
		ToSql()
	// операция идемпотентная, обновление только для открытых pr-ов

	if err != nil {
		return fmt.Errorf("build merge: %w", err)
	}

	if _, err = db.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("exec merge: %w", err)
	}
	// отсутствие обновления не является ошибкой

	return nil
}

func (r *PullRequestRepository) GetAssignedForUser(ctx context.Context, db DBTX, userID string) ([]models.PullRequest, error) {
	sql, args, err := r.builder.
		Select("pr.id", "pr.title", "pr.author_id", "pr.status").
		From("pull_requests pr").
		Join("reviewers r ON pr.id = r.pull_request_id").
		Where(squirrel.Eq{"r.reviewer_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var result []models.PullRequest
	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, pr)
	}

	return result, nil
}
