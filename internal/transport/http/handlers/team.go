package handlers

import (
	"github.com/AntonTsoy/review-pull-request-service/internal/service"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"

	"github.com/gofiber/fiber/v2"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func newTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) PostTeamAdd(c *fiber.Ctx) error {
	var req api.PostTeamAddJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Error: ErrorMessage{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	err := h.teamService.CreateTeam(c.Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(TeamResponse{Team: req})
}

func (h *TeamHandler) GetTeamGet(c *fiber.Ctx, params api.GetTeamGetParams) error {
	members, err := h.teamService.GetTeam(c.Context(), params.TeamName)
	if err != nil {
		return handleError(c, err)
	}

	resp := api.Team{
		TeamName: params.TeamName,
		Members:  members,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
