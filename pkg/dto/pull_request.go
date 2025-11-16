package dto

type PullRequest struct {
	PullRequestId   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorId        string `json:"author_id" db:"author_id"`
}

type PullRequestResponse struct {
	PullRequestId     string   `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name" db:"pull_request_name"`
	AuthorId          string   `json:"author_id" db:"author_id"`
	Status            string   `json:"status" db:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type PullRequestShort struct {
	PullRequestId   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorId        string `json:"author_id" db:"author_id"`
	Status          string `json:"status" db:"status"`
}

type UserReviewResponse struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
