package service

import (
	"errors"
	"net/http"

	"micro-service/pkg/dto"
	errs "micro-service/pkg/errors"
	rp "micro-service/pkg/repository"

	"github.com/labstack/echo/v4"
)

type PrService struct {
	rp *rp.PrRepository
}

func NewPrService(rp *rp.PrRepository) *PrService {
	return &PrService{rp: rp}
}

func (pr *PrService) CreatePullRequest(c echo.Context) error {
	var body dto.PullRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}
	pullRequest, err := pr.rp.CreatePR(body)
	if errors.Is(err, &errs.NotFoundError{}) || errors.Is(err, &errs.NoActiveReviewersError{}) {
		return c.JSON(http.StatusNotFound, err.Error())
	} else if errors.Is(err, &errs.PrExistsError{}) {
		return c.JSON(http.StatusConflict, err.Error())
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, (&errs.InternalServerError{}).Error())
	}

	return c.JSON(http.StatusCreated, pullRequest)
}

func (pr *PrService) MergePullRequest(c echo.Context) error {
	var body struct {
		PullRequestId string `json:"pull_request_id"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}
	pullRequest, err := pr.rp.MergePR(body.PullRequestId)
	if errors.Is(err, &errs.NotFoundError{}) {
		return c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, (&errs.InternalServerError{}).Error())
	}
	return c.JSON(http.StatusOK, pullRequest)
}

func (s *PrService) ReassignReviewer(c echo.Context) error {
	var req struct {
		PullRequestId string `json:"pull_request_id"`
		OldUserId     string `json:"old_user_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}

	newReviewer, err := s.rp.ReassignReviewer(req.PullRequestId, req.OldUserId)
	if errors.Is(err, &errs.InternalServerError{}) {
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else if errors.Is(err, &errs.NotFoundError{}) {
		return c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		return c.JSON(http.StatusConflict, err.Error())
	}

	resp, err := s.rp.GetPrResponse(req.PullRequestId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"pr":          resp,
		"replaced_by": newReviewer,
	})
}

func (s *PrService) GetUserReviews(c echo.Context) error {
	userId := c.QueryParam("user_id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}

	prs, err := s.rp.GetReviewsByUser(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, (&errs.InternalServerError{}).Error())
	}

	resp := dto.UserReviewResponse{
		UserId:       userId,
		PullRequests: prs,
	}

	return c.JSON(http.StatusOK, resp)
}
