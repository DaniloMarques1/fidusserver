package models

import (
	"database/sql"

	"github.com/danilomarques1/fidusserver/database"
)

type Master struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type MasterDAO interface {
	Save(*Master) error
	FindByEmail(string) (*Master, error)
}

type masterDAODatabase struct {
	db *sql.DB
}

func NewMasterDAODatabase() MasterDAO {
	db := database.Database()
	return &masterDAODatabase{db}
}

func (m *masterDAODatabase) Save(master *Master) error {
	stmt, err := m.db.Prepare(`insert into fidus_master(id, name, email, password_hash) values($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(master.ID, master.Name, master.Email, master.PasswordHash); err != nil {
		return err
	}
	return nil
}

func (m *masterDAODatabase) FindByEmail(email string) (*Master, error) {
	stmt, err := m.db.Prepare("select id, name, email, password_hash from fidus_master where email = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	master := &Master{}
	if err := stmt.QueryRow(email).Scan(&master.ID, &master.Name, &master.Email, &master.PasswordHash); err != nil {
		return nil, err
	}
	return master, nil
}
