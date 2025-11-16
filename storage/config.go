package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Config struct {
	UserName string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		UserName: viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetInt("POSTGRES_PORT"),
		DBName:   viper.GetString("POSTGRES_DB"),
		SSLMode:  viper.GetString("SSLMODE"),
	}
}

func NewStorage(c Config) (*sqlx.DB, error) {
	const op = "storage.config.NewStorage"
	db, err := sqlx.Connect(
		"pgx",
		fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
			c.UserName, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
		),
	)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу, если ее еще нет
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Users(
        user_id VARCHAR(128) PRIMARY KEY,
        username VARCHAR(256) NOT NULL,
        is_active BOOLEAN DEFAULT FALSE);

	CREATE TABLE IF NOT EXISTS Teams(
        team_name VARCHAR(256) PRIMARY KEY);

	CREATE TABLE IF NOT EXISTS TeamUsers(
        id SERIAL PRIMARY KEY,
        user_id VARCHAR(128) NOT NULL,
		team_name VARCHAR(256) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
		FOREIGN KEY (team_name) REFERENCES teams (team_name) ON DELETE CASCADE);

    CREATE TABLE IF NOT EXISTS PullRequests(
        pull_request_id VARCHAR(128) PRIMARY KEY,
        pull_request_name VARCHAR(256) NOT NULL,
        status VARCHAR(256) NOT NULL CHECK((status != '') AND ((status = 'OPEN') OR (status = 'MERGED'))),
		author_id VARCHAR(128) NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users (user_id) ON DELETE CASCADE
		);

	CREATE TABLE IF NOT EXISTS PullRequestUsers(
        id SERIAL PRIMARY KEY,
        user_id VARCHAR(128) NOT NULL,
		pull_request_id VARCHAR(128) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
		FOREIGN KEY (pull_request_id) REFERENCES PullRequests (pull_request_id) ON DELETE CASCADE);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
