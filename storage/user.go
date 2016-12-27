package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID    int
	Login string `validate:"required"`
	Pass  string `validate:"required"`
	profileInfo
}

type profileInfo struct {
	Name     string `validate:"neglect"`
	Age      int64  `validate:"neglect"`
	Phone    string `validate:"neglect"`
	CreateAt time.Time
}

const (
	sqlUserCheck  = `SELECT id, login, name, age, phone, createat FROM movies."user" WHERE login=$1 AND pass=$2`
	sqlUserValid  = `SELECT id FROM movies."user" WHERE login = $1`
	sqlUserInsert = `INSERT INTO movies."user" (login, pass, name, age, phone, createat) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
)

func (u *User) Checking() error {
	err := db.QueryRow(sqlUserCheck, u.Login, u.Pass).Scan(&u.ID, &u.Login, &u.Name, &u.Age, &u.Phone, &u.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Bad login or password")
		}
		return err
	}
	return nil
}

// Create user
func (u *User) Create() error {
	err := db.QueryRow(sqlUserValid, u.Login).Scan(&u.ID)
	if err == nil {
		return fmt.Errorf("login %v exists", u.Login)
	}

	err = db.QueryRow(sqlUserInsert, u.Login, u.Pass, u.Name, u.Age, u.Phone, u.CreateAt).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}
