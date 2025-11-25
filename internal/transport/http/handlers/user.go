package handlers

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/service"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
	prService   *service.PullRequestService
}

func newUserHandler(userService *service.UserService, prService *service.PullRequestService) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

func (h *UserHandler) GetUsersGetReview(c *fiber.Ctx, params api.GetUsersGetReviewParams) error {
	prs, err := h.prService.GetReviewForUser(c.Context(), params.UserId)
	if err != nil {
		return handleError(c, err)
	}

	shortPRs := make([]api.PullRequestShort, len(prs))
	for i, pr := range prs {
		shortPRs[i] = api.PullRequestShort{
			PullRequestId:   pr.ID,
			PullRequestName: pr.Title,
			AuthorId:        pr.AuthorID,
			Status:          api.PullRequestShortStatus(pr.Status),
		}
	}

	resp := UserGetReviewResponse{
		UserId:       params.UserId,
		PullRequests: shortPRs,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *UserHandler) PostUsersSetIsActive(c *fiber.Ctx) error {
	var req api.PostUsersSetIsActiveJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	user, err := h.userService.SetIsActive(c.Context(), req.UserId, req.IsActive)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
