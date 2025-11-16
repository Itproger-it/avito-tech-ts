package service

import "micro-service/pkg/repository"

type Service struct {
	Team        *TeamService
	PullRequest *PrService
}

func NewService(rp *repository.Repository) *Service {
	return &Service{
		Team:        NewTeamService(rp.Team),
		PullRequest: NewPrService(rp.PullRequest),
	}
}
