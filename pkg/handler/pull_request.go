package handler

import "github.com/labstack/echo/v4"

func (h *Handler) CreatePR(c echo.Context) error {
	return h.services.PullRequest.CreatePullRequest(c)
}

func (h *Handler) MergePR(c echo.Context) error {
	return h.services.PullRequest.MergePullRequest(c)
}

func (h *Handler) GetReviewsByUser(c echo.Context) error {
	return h.services.PullRequest.GetUserReviews(c)
}
