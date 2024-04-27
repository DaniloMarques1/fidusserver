package models

import (
	"database/sql"

	"github.com/danilomarques1/fidusserver/database"
)

type Password struct {
	MasterId      string
	Key           string
	PasswordValue string
}

type PasswordDAO interface {
	Save(*Password) error
	FindOne(masterId, key string) (*Password, error)
}

type passwordDAODatabase struct {
	db *sql.DB
}

func NewPasswordDAODatabase() PasswordDAO {
	db := database.Database()
	return &passwordDAODatabase{db}
}

func (passwordDAO *passwordDAODatabase) Save(password *Password) error {
	stmt, err := passwordDAO.db.Prepare(`insert into fidus_password(master_id, key, password) values($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(password.MasterId, password.Key, password.PasswordValue); err != nil {
		return err
	}
	return nil
}

func (passwordDAO *passwordDAODatabase) FindOne(masterid, key string) (*Password, error) {
	stmt, err := passwordDAO.db.Prepare(`select master_id, key, password from fidus_password where master_id = $1 and key = $2ao`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	password := &Password{}
	if err := stmt.QueryRow().Scan(&password.MasterId, &password.Key, &password.PasswordValue); err != nil {
		return nil, err
	}
	return password, nil
}
