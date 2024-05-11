package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

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
	Delete(masterId, key string) error
	UpdatePasswordValue(masterId, key, passwordValue string) error

	// NoMatchError returns true if the error received was because it could not find a match
	NoMatchError(err error) bool
}

type passwordDAODatabase struct {
	db                *sql.DB
	passwordSecretKey string
}

func NewPasswordDAODatabase() PasswordDAO {
	db := database.Database()
	passwordSecretKey := os.Getenv("PASSWORD_ENCRYPT_KEY")
	return &passwordDAODatabase{db, passwordSecretKey}
}

func (passwordDAO *passwordDAODatabase) Save(password *Password) error {
	q := fmt.Sprintf(
		`insert into fidus_password(master_id, key, password)
		values($1, $2, pgp_sym_encrypt($3, '%s'))
		`, passwordDAO.passwordSecretKey)
	stmt, err := passwordDAO.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(password.MasterId, password.Key, password.PasswordValue); err != nil {
		return err
	}
	return nil
}

func (passwordDAO *passwordDAODatabase) FindOne(masterId, key string) (*Password, error) {
	q := fmt.Sprintf(`
	select master_id, key, pgp_sym_decrypt(password, '%s') as password
	from fidus_password where master_id = $1 and key = $2
	`, passwordDAO.passwordSecretKey)
	stmt, err := passwordDAO.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	password := &Password{}
	if err := stmt.QueryRow(masterId, key).Scan(&password.MasterId, &password.Key, &password.PasswordValue); err != nil {
		return nil, err
	}
	return password, nil
}

func (passwordDAO *passwordDAODatabase) Delete(masterId, key string) error {
	stmt, err := passwordDAO.db.Prepare(`delete from fidus_password where master_id = $1 and key = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(masterId, key); err != nil {
		return err
	}
	return nil
}

func (passwordDAO *passwordDAODatabase) UpdatePasswordValue(masterId, key, passwordValue string) error {
	q := fmt.Sprintf(`
	update fidus_password set password=pgp_sym_encrypt($1, '%s') where master_id = $2 and key = $3
	`, passwordDAO.passwordSecretKey)
	stmt, err := passwordDAO.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(passwordValue, masterId, key); err != nil {
		return err
	}
	return nil
}

func (passwordDAO *passwordDAODatabase) NoMatchError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
