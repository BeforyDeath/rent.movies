package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/BeforyDeath/rent.movies/pagination"
	"github.com/lib/pq"
)

type Rent struct {
	Id       int
	User_id  int
	Movie_id int64 `validate:"required"`
	Active   bool
	CreateAt time.Time
	CloseAt  time.Time
}

type RentList struct {
	Active   bool
	CreateAt time.Time
	CloseAt  time.Time
	Movie_id int

	Name        string
	Year        int
	Genre       string
	Description string
}

const (
	SQL_RENT_VALID  = `SELECT createAt FROM movies.rent WHERE user_id = $1 AND movie_id = $2 AND active = $3`
	SQL_RENT_INSERT = `INSERT INTO movies.rent (user_id, movie_id, active, createAt) VALUES ($1, $2, $3, $4) RETURNING id`
	SQL_RENT_UPDATE = `UPDATE movies.rent SET active=$3, closeAt=$4 WHERE active = true AND user_id = $1 AND movie_id = $2`

	SQL_RENT_SELECT = `SELECT r.active, r.createAt, r.closeAt, m.id as movie_id, m.name, m.year,
						array_to_string(movies.array_accum(g.name), ', ') AS genre, m.description `
	SQL_RENT_SELECT_COUNT = `SELECT COUNT(DISTINCT (r.id)) `
	SQL_RENT_FROM         = `FROM movies.rent r, movies.movie m, movies.movie_genre mg, movies.genre g
						WHERE r.movie_id = m.id AND mg.movie_id = m.id AND mg.genre_id = g.id
						AND r.user_id = $1 AND r.active = $2 `
	SQL_RENT_GROUP = `GROUP BY r.id, m.id `
	SQL_RENT_LIMIT = `LIMIT $3 OFFSET $4`
)

func (r Rent) Take() error {

	var CreateAt time.Time
	err := db.QueryRow(SQL_RENT_VALID, r.User_id, r.Movie_id, true).Scan(&CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {

			tx, err := db.Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			err = tx.QueryRow(SQL_RENT_INSERT, r.User_id, r.Movie_id, true, r.CreateAt).Scan(&r.Id)
			if err != nil {
				if strings.HasSuffix(err.Error(), "\"rent_movie_id_fk\"") {
					return errors.New(fmt.Sprintf("An identifier of the movie %v does not exist", r.Movie_id))
				}
				return err
			}

			err = tx.Commit()
			if err != nil {
				return err
			}

			return nil

		} else {
			return err
		}
	}
	return errors.New(fmt.Sprintf("This movie you've already rented %v", CreateAt.Format("02-01-2006 15:04")))
}

func (r Rent) GetTotalCount(user_id int, active bool) (tc int, err error) {
	err = db.QueryRow(SQL_RENT_SELECT_COUNT+SQL_RENT_FROM, user_id, active).Scan(&tc)
	return
}

func (r Rent) GetAll(p *pagination.Pages, user_id int, active bool) ([]RentList, error) {
	rows, err := db.Query(SQL_RENT_SELECT+SQL_RENT_FROM+SQL_RENT_GROUP+SQL_RENT_LIMIT, user_id, active, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]RentList, 0)
	for rows.Next() {

		var CloseAt pq.NullTime

		rl := RentList{}
		err := rows.Scan(&rl.Active, &rl.CreateAt, &CloseAt, &rl.Movie_id, &rl.Name, &rl.Year, &rl.Genre, &rl.Description)
		if err != nil {
			return nil, err
		}

		if CloseAt.Valid {
			rl.CloseAt = CloseAt.Time
		}

		res = append(res, rl)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Rent) Completed() error {
	err := db.QueryRow(SQL_RENT_VALID, r.User_id, r.Movie_id, true).Scan(&r.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("An identifier of the movie %v is not leased", r.Movie_id))
		}
		return err
	}

	res, err := db.Exec(SQL_RENT_UPDATE, r.User_id, r.Movie_id, false, time.Now())
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
