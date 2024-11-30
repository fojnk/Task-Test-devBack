package repository

import (
	"fmt"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user models.User) (string, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, email, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := a.db.QueryRow(query, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return string(id), nil
}

func (a *AuthPostgres) GetUser(guid string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE guid=$1", usersTable)

	err := a.db.Get(&user, query, guid)
	return user, err
}
