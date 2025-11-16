package handler

import "github.com/labstack/echo/v4"

func (h *Handler) TeadAdd(c echo.Context) error {
	return h.services.Team.CreateTeam(c)
}

func (h *Handler) TeadGet(c echo.Context) error {
	return h.services.Team.GetTeam(c)
}

func (h *Handler) UserSetIsActive(c echo.Context) error {
	return h.services.Team.TeamUserIsActive(c)
}
