package apperrors

import "errors"

var (
	ErrNotFound   = errors.New("resource not found")
	ErrTeamExists = errors.New("team_name already exists")
)
