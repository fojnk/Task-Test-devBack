package repository

import (
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GetUser(guid string) (models.User, error)
}

type Respository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Respository {
	return &Respository{
		Authorization: NewAuthPostgres(db),
	}
}
