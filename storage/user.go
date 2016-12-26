package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id       int
	Login    string `validate:"required"`
	Pass     string `validate:"required"`
	Name     string `validate:"neglect"`
	Age      int64  `validate:"neglect"`
	Phone    string `validate:"neglect"`
	CreateAt time.Time
}

const (
	SQL_USER_CHECK  = `SELECT id, login, name, age, phone, createat FROM movies."user" WHERE login=$1 AND pass=$2`
	SQL_USER_VALID  = `SELECT id FROM movies."user" WHERE login = $1`
	SQL_USER_INSERT = `INSERT INTO movies."user" (login, pass, name, age, phone, createat) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
)

func (u *User) Check() error {
	err := db.QueryRow(SQL_USER_CHECK, u.Login, u.Pass).Scan(&u.Id, &u.Login, &u.Name, &u.Age, &u.Phone, &u.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Bad login or password")
		}
		return err
	}
	return nil
}

func (u *User) Create() error {
	err := db.QueryRow(SQL_USER_VALID, u.Login).Scan(&u.Id)
	if err == nil {
		return errors.New(fmt.Sprintf("login %v exists", u.Login))
	}

	err = db.QueryRow(SQL_USER_INSERT, u.Login, u.Pass, u.Name, u.Age, u.Phone, u.CreateAt).Scan(&u.Id)
	if err != nil {
		return err
	}
	return nil
}
