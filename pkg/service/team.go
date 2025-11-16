package service

import (
	"errors"
	"net/http"

	"micro-service/pkg/dto"
	errs "micro-service/pkg/errors"
	rp "micro-service/pkg/repository"

	"github.com/labstack/echo/v4"
)

type TeamService struct {
	rp *rp.TeamRepository
}

func NewTeamService(rp *rp.TeamRepository) *TeamService {
	return &TeamService{rp: rp}
}

func (t *TeamService) CreateTeam(c echo.Context) error {
	team := new(dto.Team)
	if err := c.Bind(team); err != nil {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}
	err := t.rp.CreateTeam(*team)
	if errors.Is(err, &errs.TeamExistsError{}) {
		return c.JSON(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, (&errs.InternalServerError{}).Error())
	}
	return c.JSON(http.StatusCreated, team)
}

func (t *TeamService) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}
	teamMembers, err := t.rp.TeamUsers(teamName)
	if errors.Is(err, &errs.NotFoundError{}) {
		return c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, (&errs.InternalServerError{}).Error())
	}
	return c.JSON(http.StatusOK, teamMembers)
}

func (t *TeamService) TeamUserIsActive(c echo.Context) error {
	userIsActive := new(dto.TeamMemberIsActive)
	if err := c.Bind(userIsActive); err != nil {
		return c.JSON(http.StatusBadRequest, (&errs.BadRequestError{}).Error())
	}
	teamMember, err := t.rp.TeamUserIsActive(*userIsActive)
	if errors.Is(err, &errs.NotFoundError{}) {
		return c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, (&errs.InternalServerError{}).Error())
	}
	return c.JSON(http.StatusOK, teamMember)
}
