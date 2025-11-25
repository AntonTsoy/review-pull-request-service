package handlers

import "github.com/AntonTsoy/review-pull-request-service/internal/service"

type TeamHanler struct {
	teamService *service.TeamService
}

func newTeamHandler(teamService *service.TeamService) *TeamHanler {
	return &TeamHanler{
		teamService: teamService,
	}
}

