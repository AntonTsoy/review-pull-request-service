package models

import "time"

type StatusPR = string

const (
	StatusOpen   StatusPR = "OPEN"
	StatusMerged StatusPR = "MERGED"
)

type PullRequest struct {
	ID                string
	Title             string
	AuthorID          string
	Status            StatusPR
	MergedAt          *time.Time
	AssignedReviewers []string
}
