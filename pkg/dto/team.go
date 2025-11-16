package dto

type TeamMember struct {
	UserId   string `json:"user_id" db:"user_id" binding:"required"`
	UserName string `json:"username" db:"username" binding:"required"`
	IsActive bool   `json:"is_active" db:"is_active" binding:"required"`
}

type TeamMemberIsActive struct {
	UserId   string `json:"user_id" db:"user_id" binding:"required"`
	IsActive bool   `json:"is_active" db:"is_active" binding:"required"`
}

type TeamAndUserID struct {
	UserId   string `json:"user_id" db:"user_id" binding:"required"`
	TeamName string `json:"team_name" db:"team_name" binding:"required"`
}

type Team struct {
	TeamName string       `json:"team_name" db:"team_name" binding:"required"`
	Members  []TeamMember `json:"members" db:"members" binding:"required"`
}

type TeamAndMember struct {
	TeamName string `json:"team_name" db:"-" binding:"required"`
	UserId   string `json:"user_id" db:"user_id" binding:"required"`
	UserName string `json:"username" db:"username" binding:"required"`
	IsActive bool   `json:"is_active" db:"is_active" binding:"required"`
}

type TeamName struct {
	TeamName string `json:"team_name" query:"team_name" db:"team_name" binding:"required"`
}
