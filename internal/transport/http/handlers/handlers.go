package handlers

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/apperrors"
	"github.com/AntonTsoy/review-pull-request-service/internal/service"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"

	"github.com/gofiber/fiber/v2"
)

type handlers struct {
	*TeamHandler
	*UserHandler
	*PullRequestHandler
}

func NewHandlers(service *service.Service) api.ServerInterface {
	return &handlers{
		TeamHandler:        newTeamHandler(service.TeamService),
		UserHandler:        newUserHandler(service.UserService, service.PullRequestService),
		PullRequestHandler: newPullRequestHandler(service.PullRequestService),
	}
}

func handleError(c *fiber.Ctx, err error) error {
	switch err {
	case apperrors.ErrTeamExists:
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "TEAM_EXISTS",
				Message: err.Error(),
			},
		})
	case apperrors.ErrPullRequestExists:
		return c.Status(fiber.StatusConflict).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "PR_EXISTS",
				Message: err.Error(),
			},
		})
	case apperrors.ErrPullRequestMerged:
		return c.Status(fiber.StatusConflict).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "PR_MERGED",
				Message: err.Error(),
			},
		})
	case apperrors.ErrNotAssigned:
		return c.Status(fiber.StatusConflict).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "NOT_ASSIGNED",
				Message: err.Error(),
			},
		})
	case apperrors.ErrNoCandidate:
		return c.Status(fiber.StatusConflict).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "NO_CANDIDATE",
				Message: err.Error(),
			},
		})
	case apperrors.ErrNotFound:
		return c.Status(fiber.StatusNotFound).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
	}
}
