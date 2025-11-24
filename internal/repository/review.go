package repository

import (
	"context"
	"fmt"

	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/Masterminds/squirrel"
)

type ReviewRepository struct {
	builder squirrel.StatementBuilderType
}

func newReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ReviewRepository) Assign(ctx context.Context, db DBTX, prID string, reviewerIDs ...string) error {
	query := r.builder.
		Insert("reviewers").
		Columns("name", "email", "age")
	for _, reviewerID := range reviewerIDs {
		query = query.Values(prID, reviewerID)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("insert new reviewer: %w", err)
	}

	return nil
}

func (r *ReviewRepository) GetReviewersByPR(ctx context.Context, db DBTX, prID string) ([]string, error) {
	sql, args, err := r.builder.
		Select("reviewer_id").
		From("reviewers").
		Where(squirrel.Eq{"pull_request_id": prID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build reviewers query: %w", err)
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query reviewers: %w", err)
	}
	defer rows.Close()

	var reviewerIDs []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		reviewerIDs = append(reviewerIDs, reviewerID)
	}

	return reviewerIDs, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, db DBTX, prID, oldReviewerID string) error {
	sql, args, err := r.builder.
		Delete("reviewers").
		Where(squirrel.Eq{"pull_request_id": prID, "reviewer_id": oldReviewerID}).
		ToSql()
	
	if err != nil {
		return fmt.Errorf("build delete: %w", err)
	}

	tag, err := db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete old reviewer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
