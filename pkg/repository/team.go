package repository

import (
	"micro-service/pkg/dto"
	errs "micro-service/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeam(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (t *TeamRepository) CreateTeam(team dto.Team) error {
	if err := t.GetTeam(team.TeamName); err == nil {
		return &errs.TeamExistsError{}
	}
	for _, team_member := range team.Members {
		_, err := t.GetUser(team_member.UserId)
		if err != nil {
			t.AddUser(team_member)
		}
	}

	query := "INSERT INTO Teams (team_name) values($1);"
	_, err := t.db.Exec(query, team.TeamName)
	if err != nil {
		return err
	}

	for _, teamMember := range team.Members {
		t.UpdateUser(team.TeamName, teamMember)
	}

	return nil
}

func (t *TeamRepository) GetTeam(team_name string) error {
	var team dto.TeamName
	query := "SELECT team_name FROM Teams WHERE team_name = $1;"
	err := t.db.Get(&team, query, team_name)
	if err != nil {
		return &errs.NotFoundError{}
	}
	return nil
}

func (t *TeamRepository) TeamUsers(team_name string) ([]dto.TeamMember, error) {
	err := t.GetTeam(team_name)
	if err != nil {
		return nil, &errs.NotFoundError{}
	}

	var users []dto.TeamMember
	query := `
		SELECT u.user_id, u.username, u.is_active
		FROM Users u
		JOIN TeamUsers tu ON u.user_id = tu.user_id
		WHERE tu.team_name = $1;
	`

	err = t.db.Select(&users, query, team_name)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (t *TeamRepository) AddUser(team_member dto.TeamMember) error {
	query := "INSERT INTO Users (user_id, username, is_active) values($1, $2, $3);"
	_, err := t.db.Exec(query, team_member.UserId, team_member.UserName, team_member.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (t *TeamRepository) GetUser(user_id string) (*dto.TeamMember, error) {
	var get_user = new(dto.TeamMember)
	query := "SELECT * FROM Users WHERE user_id = $1;"
	err := t.db.Get(get_user, query, user_id)
	return get_user, err
}

func (t *TeamRepository) UpdateUser(teamName string, tm dto.TeamMember) error {
	// обновление пользователя
	_, err := t.db.Exec(
		"UPDATE Users SET username = $1, is_active = $2 WHERE user_id = $3;",
		tm.UserName, tm.IsActive, tm.UserId,
	)
	if err != nil {
		return err
	}

	// Проверяем связь user → team
	var bind dto.TeamAndUserID
	err = t.db.Get(&bind, "SELECT user_id, team_name FROM TeamUsers WHERE user_id = $1;", tm.UserId)
	if err != nil { // связи нет — создаём
		_, err = t.db.Exec(
			"INSERT INTO TeamUsers (user_id, team_name) VALUES ($1, $2);",
			tm.UserId, teamName,
		)
		return err
	}

	// связь есть — обновляем команду
	_, err = t.db.Exec(
		"UPDATE TeamUsers SET team_name = $1 WHERE user_id = $2;",
		teamName, tm.UserId,
	)
	return err
}

func (t *TeamRepository) TeamUserIsActive(data dto.TeamMemberIsActive) (*dto.TeamAndMember, error) {

	_, err := t.db.Exec(
		"UPDATE Users SET is_active = $2 WHERE user_id = $1;",
		data.UserId, data.IsActive,
	)
	if err != nil {
		return nil, err
	}

	member, err := t.GetUser(data.UserId)
	if err != nil {
		return nil, &errs.NotFoundError{}
	}

	var team dto.TeamName
	query := `
		SELECT t.team_name FROM Teams t
		JOIN TeamUsers tu ON t.team_name = tu.team_name
		WHERE tu.user_id = $1;
	`
	err = t.db.Get(&team, query, member.UserId)
	if err != nil {
		return nil, &errs.NotFoundError{}
	}
	return &dto.TeamAndMember{
		TeamName: team.TeamName,
		UserId:   member.UserId,
		UserName: member.UserName,
		IsActive: member.IsActive,
	}, nil
}
