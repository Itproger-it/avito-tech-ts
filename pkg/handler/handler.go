package handler

import (
	"micro-service/pkg/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) InitRoutes() *echo.Echo {
	e := echo.New()

	e.POST("/team/add", func(c echo.Context) error {
		return h.services.Team.CreateTeam(c)
	})

	e.GET("/team/get", func(c echo.Context) error {
		return h.services.Team.GetTeam(c)
	})

	e.POST("/users/setIsActive", func(c echo.Context) error {
		return h.services.Team.TeamUserIsActive(c)
	})

	e.POST("/pullRequest/create", func(c echo.Context) error {
		return h.services.PullRequest.CreatePullRequest(c)
	})

	e.POST("/pullRequest/merge", func(c echo.Context) error {
		return h.services.PullRequest.MergePullRequest(c)
	})

	e.POST("/pullRequest/reassign", func(c echo.Context) error {
		return h.services.PullRequest.ReassignReviewer(c)
	})

	e.GET("/users/getReview", func(c echo.Context) error {
		return h.services.PullRequest.GetUserReviews(c)
	})

	return e
}
