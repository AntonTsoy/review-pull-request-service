package handlers

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/service"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"

	"github.com/gofiber/fiber/v2"
)

type PullRequestHandler struct {
	prService *service.PullRequestService
}

func newPullRequestHandler(prService *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		prService: prService,
	}
}

func (h *PullRequestHandler) PostPullRequestCreate(c *fiber.Ctx) error {
	var req api.PostPullRequestCreateJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	pr, err := h.prService.Create(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(PullRequestResponse{Pr: convertPRToAPI(pr)})
}

func (h *PullRequestHandler) PostPullRequestMerge(c *fiber.Ctx) error {
	var req api.PostPullRequestMergeJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	pr, err := h.prService.Merge(c.Context(), req.PullRequestId)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(PullRequestResponse{Pr: convertPRToAPI(pr)})
}

func (h *PullRequestHandler) PostPullRequestReassign(c *fiber.Ctx) error {
	var req api.PostPullRequestReassignJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	pr, replacedBy, err := h.prService.Reassign(c.Context(), req.PullRequestId, req.OldUserId)
	if err != nil {
		return handleError(c, err)
	}

	resp := PullRequestAssignResponse{
		Pr:         convertPRToAPI(pr),
		ReplacedBy: replacedBy,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
