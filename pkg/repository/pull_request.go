package repository

import (
	dto "micro-service/pkg/dto"
	errs "micro-service/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type PrRepository struct {
	db *sqlx.DB
}

func NewPR(db *sqlx.DB) *PrRepository {
	return &PrRepository{db: db}
}

func (pr *PrRepository) CreatePR(data dto.PullRequest) (*dto.PullRequestResponse, error) {
	tx, err := pr.db.Beginx()
	if err != nil {
		return nil, err
	}

	var exists string
	err = tx.Get(&exists, "SELECT pull_request_id FROM PullRequests WHERE pull_request_id = $1;", data.PullRequestId)
	if err == nil {
		tx.Rollback()
		return nil, &errs.PrExistsError{}
	}

	var teamName string
	err = tx.Get(&teamName,
		`SELECT team_name FROM TeamUsers WHERE user_id = $1 LIMIT 1;`,
		data.AuthorId,
	)
	if err != nil {
		tx.Rollback()
		return nil, &errs.NotFoundError{}
	}

	reviewers := []string{}
	err = tx.Select(&reviewers,
		`SELECT user_id FROM Users 
         WHERE user_id IN (
             SELECT user_id FROM TeamUsers WHERE team_name = $1
         )
         AND is_active = true
         AND user_id != $2
         LIMIT 2;`,
		teamName, data.AuthorId,
	)
	if err != nil || len(reviewers) == 0 {
		tx.Rollback()
		return nil, &errs.NoActiveReviewersError{}
	}

	_, err = tx.Exec(
		`INSERT INTO PullRequests (pull_request_id, pull_request_name, author_id, status)
         VALUES ($1, $2, $3, 'OPEN');`,
		data.PullRequestId, data.PullRequestName, data.AuthorId,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, r := range reviewers {
		_, err = tx.Exec(
			`INSERT INTO PullRequestUsers (pull_request_id, user_id) VALUES ($1, $2);`,
			data.PullRequestId, r,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return &dto.PullRequestResponse{
		PullRequestId:     data.PullRequestId,
		PullRequestName:   data.PullRequestName,
		AuthorId:          data.AuthorId,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
	}, nil
}

func (pr *PrRepository) MergePR(id string) (*dto.PullRequestResponse, error) {
	var prData dto.PullRequestResponse
	err := pr.db.Get(&prData,
		`SELECT pull_request_id, pull_request_name, author_id, status 
         FROM PullRequests WHERE pull_request_id = $1;`,
		id,
	)
	if err != nil {
		return nil, &errs.NotFoundError{}
	}

	if prData.Status == "MERGED" {
		return pr.loadReviewers(prData)
	}

	_, err = pr.db.Exec(
		`UPDATE PullRequests SET status='MERGED' WHERE pull_request_id=$1;`,
		id,
	)
	if err != nil {
		return nil, err
	}

	prData.Status = "MERGED"
	return pr.loadReviewers(prData)
}

func (pr *PrRepository) ReassignReviewer(prID string, oldUser string) (string, error) {
	var status string
	err := pr.db.Get(&status,
		"SELECT status FROM PullRequests WHERE pull_request_id = $1",
		prID,
	)
	if err != nil {
		return "", &errs.NotFoundError{}
	}
	if status == "MERGED" {
		return "", &errs.PrMergedError{}
	}

	var exists bool
	err = pr.db.Get(&exists,
		"SELECT EXISTS(SELECT 1 FROM PullRequestUsers WHERE pull_request_id=$1 AND user_id=$2)",
		prID, oldUser,
	)
	if err != nil || !exists {
		return "", &errs.NotAssignedError{}
	}

	var teamAndAuthor struct {
		TeamName string `json:"team_name" db:"team_name"`
		UserId   string `json:"user_id" db:"user_id"`
	}
	err = pr.db.Get(&teamAndAuthor,
		`SELECT tu.team_name, tu.user_id
		 FROM PullRequests pr
		 JOIN TeamUsers tu ON tu.user_id = pr.author_id
		 WHERE pr.pull_request_id=$1`,
		prID,
	)
	if err != nil {
		return "", &errs.NotFoundError{}
	}

	candidates := []string{}
	err = pr.db.Select(&candidates,
		`SELECT u.user_id
		 FROM Users u
		 JOIN TeamUsers tu ON tu.user_id = u.user_id
		 WHERE tu.team_name=$1
		   AND u.is_active = TRUE
		   AND u.user_id NOT IN ($2, $4)
		   AND u.user_id NOT IN (
				SELECT user_id FROM PullRequestUsers WHERE pull_request_id=$3
		   )`,
		teamAndAuthor.TeamName, oldUser, prID, teamAndAuthor.UserId,
	)
	if err != nil || len(candidates) == 0 {
		return "", &errs.NoCandidateError{}
	}

	newReviewer := candidates[0] // можно сделать random

	_, err = pr.db.Exec(
		`UPDATE PullRequestUsers
		 SET user_id = $1
		 WHERE pull_request_id=$2 AND user_id=$3`,
		newReviewer, prID, oldUser,
	)
	if err != nil {
		return "", &errs.InternalServerError{}
	}
	return newReviewer, nil
}

func (pr *PrRepository) GetPrResponse(prID string) (*dto.PullRequestResponse, error) {
	var prInfo dto.PullRequestResponse

	err := pr.db.Get(&prInfo,
		`SELECT pull_request_id, pull_request_name, author_id, status
		 FROM PullRequests WHERE pull_request_id=$1`,
		prID,
	)
	if err != nil {
		return nil, err
	}

	err = pr.db.Select(&prInfo.AssignedReviewers,
		`SELECT user_id FROM PullRequestUsers WHERE pull_request_id=$1`,
		prID,
	)
	if err != nil {
		return nil, err
	}

	return &prInfo, nil
}

func (pr *PrRepository) loadReviewers(prData dto.PullRequestResponse) (*dto.PullRequestResponse, error) {
	reviewers := []string{}
	err := pr.db.Select(&reviewers,
		`SELECT user_id FROM PullRequestUsers WHERE pull_request_id = $1;`,
		prData.PullRequestId,
	)
	if err != nil {
		return nil, err
	}
	prData.AssignedReviewers = reviewers
	return &prData, nil
}

func (pr *PrRepository) GetUserReviews(userId string) ([]dto.PullRequestResponse, error) {
	prs := []dto.PullRequestResponse{}
	err := pr.db.Select(&prs,
		`SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
         FROM PullRequests pr
         JOIN PullRequestUsers pu ON pr.pull_request_id = pu.pull_request_id
         WHERE pu.user_id = $1;`,
		userId,
	)
	return prs, err
}

func (pr *PrRepository) GetReviewsByUser(userId string) ([]dto.PullRequestShort, error) {
	prs := []dto.PullRequestShort{}

	query := `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM PullRequests pr
		JOIN PullRequestUsers pu ON pr.pull_request_id = pu.pull_request_id
		WHERE pu.user_id = $1
		ORDER BY pr.pull_request_id;
	`

	err := pr.db.Select(&prs, query, userId)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
