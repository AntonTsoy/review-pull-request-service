package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/database"
	"github.com/AntonTsoy/review-pull-request-service/internal/models"
	"github.com/AntonTsoy/review-pull-request-service/internal/repository"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
	"github.com/jackc/pgx/v5"
)

type PullRequestService struct {
	db         *database.Database
	prRepo     *repository.PullRequestRepository
	userRepo   *repository.UserRepository
	reviewRepo *repository.ReviewRepository
}

func newPullRequestService(
	db *database.Database,
	prRepo *repository.PullRequestRepository,
	userRepo *repository.UserRepository,
	reviewRepo *repository.ReviewRepository,
) *PullRequestService {
	return &PullRequestService{
		db:         db,
		prRepo:     prRepo,
		userRepo:   userRepo,
		reviewRepo: reviewRepo,
	}
}

func (s *PullRequestService) Create(ctx context.Context, req *api.PostPullRequestCreateJSONRequestBody) (*models.PullRequest, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	existsPR, err := s.prRepo.Exists(ctx, tx, req.PullRequestId)
	if err == nil && existsPR {
		return nil, apperrors.ErrPullRequestExists
	}

	existsAuthor, err := s.userRepo.Exists(ctx, tx, req.AuthorId)
	if err == nil && !existsAuthor {
		return nil, apperrors.ErrNotFound
	}

	candidates, err := s.userRepo.GetActiveTeammates(ctx, tx, req.AuthorId)
	if err != nil {
		return nil, err
	}

	pr := &models.PullRequest{
		ID:                req.PullRequestId,
		Title:             req.PullRequestName,
		AuthorID:          req.AuthorId,
		Status:            models.StatusOpen,
		AssignedReviewers: chooseRandomReviewers(candidates),
	}

	if err = s.prRepo.Create(ctx, tx, pr); err != nil {
		return nil, err
	}

	if err = s.reviewRepo.Assign(ctx, tx, req.PullRequestId, pr.AssignedReviewers...); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return pr, nil
}

func (s *PullRequestService) Merge(ctx context.Context, prID string) (*models.PullRequest, error) {
	if err := s.prRepo.UpdateMergeStatus(ctx, s.db.Pool(), prID); err != nil {
		return nil, err
	}

	pr, err := s.prRepo.GetByID(ctx, s.db.Pool(), prID)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers, err = s.reviewRepo.GetReviewersByPR(ctx, s.db.Pool(), prID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PullRequestService) Reassign(ctx context.Context, prID, oldUserID string) (*models.PullRequest, string, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	pr, err := s.prRepo.GetByID(ctx, tx, prID)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == models.StatusMerged {
		return nil, "", apperrors.ErrPullRequestMerged
	}

	if pr.AssignedReviewers, err = s.reviewRepo.GetReviewersByPR(ctx, tx, prID); err != nil {
		return nil, "", err
	}

	if err = checkAssignedUser(pr.AssignedReviewers, oldUserID); err != nil {
		return nil, "", err
	}

	candidates, err := s.userRepo.GetActiveTeammates(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, "", err
	}

	newReviewer, err := changeAvailableReviewer(candidates, pr.AssignedReviewers)
	if err != nil {
		return nil, "", err
	}

	if err = s.reviewRepo.Delete(ctx, tx, prID, oldUserID); err != nil {
		return nil, "", err
	}

	if err = s.reviewRepo.Assign(ctx, tx, prID, newReviewer); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	updatedPR, _ := s.prRepo.GetByID(ctx, s.db.Pool(), prID)
	return updatedPR, newReviewer, nil
}

func (s *PullRequestService) GetReviewForUser(ctx context.Context, userID string) ([]models.PullRequest, error) {
	return s.prRepo.GetAssignedForUser(ctx, s.db.Pool(), userID)
}

func chooseRandomReviewers(candidates []api.TeamMember) []string {
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	reviewers := make([]string, min(2, len(candidates)))
	for i := 0; i < len(reviewers); i++ {
		reviewers[i] = candidates[i].UserId
	}

	return reviewers
}

func checkAssignedUser(assignedReviewers []string, oldUserID string) error {
	assigned := false
	for _, r := range assignedReviewers {
		if r == oldUserID {
			assigned = true
			break
		}
	}
	if !assigned {
		return apperrors.ErrNotAssigned
	}

	return nil
}

func changeAvailableReviewer(candidates []api.TeamMember, assignedReviewers []string) (string, error) {
	available := make(map[string]struct{}, len(candidates) + 2)
	for _, candidate := range candidates {
		available[candidate.UserId] = struct{}{}
	}
	for _, currReviewer := range assignedReviewers {
		delete(available, currReviewer)
	}

	if len(available) == 0 {
		return "", apperrors.ErrNoCandidate
	}

	for newReviewer := range available { // ха-ха-ха, я пользуюсь тем, что в Go случайный обход мапы
		return newReviewer, nil
	}

	return "", nil // заглушка, чтобы gopls не ругался
}
