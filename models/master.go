package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/danilomarques1/fidusserver/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Master struct {
	ID                     string
	Name                   string
	Email                  string
	PasswordHash           string
	CreatedAt              time.Time
	PasswordExpirationDate time.Time
}

func NewMaster(name, email, password string) (*Master, error) {
	master := &Master{Name: name, Email: email}
	if err := master.HashPassword(password); err != nil {
		return nil, err
	}
	master.SetPasswordExpiration()
	master.ID = uuid.NewString()
	return master, nil
}

func (m *Master) SetPasswordExpiration() {
	m.PasswordExpirationDate = time.Now().Add(2190 * time.Hour)
}

func (m *Master) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(m.PasswordHash), []byte(password))
}

func (m *Master) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.PasswordHash = string(hashed)
	return nil
}

func (m *Master) IsPasswordExpired() bool {
	now := time.Now()
	return now.After(m.PasswordExpirationDate)
}

type MasterDAO interface {
	Save(*Master) error
	FindByEmail(string) (*Master, error)
	FindById(string) (*Master, error)
	ResetPassword(masterId, password string, passwordExpirationDate time.Time) error

	// NoMatchError returns true if the error received was because it could find a match
	NoMatchError(err error) bool
}

type masterDAODatabase struct {
	db *sql.DB
}

func NewMasterDAODatabase() MasterDAO {
	db := database.Database()
	return &masterDAODatabase{db}
}

func (m *masterDAODatabase) Save(master *Master) error {
	stmt, err := m.db.Prepare(`insert into fidus_master(id, name, email, password_hash, password_expiration_date) values($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(master.ID, master.Name, master.Email, master.PasswordHash, master.PasswordExpirationDate); err != nil {
		return err
	}
	return nil
}

func (m *masterDAODatabase) FindByEmail(email string) (*Master, error) {
	stmt, err := m.db.Prepare("select id, name, email, password_hash, password_expiration_date from fidus_master where email = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	master := &Master{}
	if err := stmt.QueryRow(email).Scan(&master.ID, &master.Name, &master.Email, &master.PasswordHash, &master.PasswordExpirationDate); err != nil {
		return nil, err
	}
	return master, nil
}

func (m *masterDAODatabase) FindById(masterId string) (*Master, error) {
	stmt, err := m.db.Prepare("select id, name, email, password_hash, password_expiration_date from fidus_master where id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	master := &Master{}
	if err := stmt.QueryRow(masterId).Scan(&master.ID, &master.Name, &master.Email, &master.PasswordHash, &master.PasswordExpirationDate); err != nil {
		return nil, err
	}
	return master, nil
}

func (m *masterDAODatabase) ResetPassword(masterId, password string, passwordExpirationDate time.Time) error {
	stmt, err := m.db.Prepare(`update fidus_master set password_hash = $1, password_expiration_date = $2 where id = $3`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(password, passwordExpirationDate, masterId); err != nil {
		return err
	}
	return nil
}

func (m *masterDAODatabase) NoMatchError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
