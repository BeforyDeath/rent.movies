package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/BeforyDeath/rent.movies/pagination"
	"github.com/lib/pq"
)

type createCloseAt struct {
	CreateAt time.Time
	CloseAt  time.Time
}

type Rent struct {
	ID      int
	UserID  int   `json:"-"`
	MovieID int64 `validate:"required"`
	Active  bool
	createCloseAt
}

type RentMovie struct {
	Rent
	Movie
}

const (
	sqlRentValid  = `SELECT createAt FROM movies.rent WHERE userId = $1 AND movieId = $2 AND active = $3`
	sqlRentInsert = `INSERT INTO movies.rent (userId, movieId, active, createAt) VALUES ($1, $2, $3, $4) RETURNING id`
	sqlRentUpdate = `UPDATE movies.rent SET active=$3, closeAt=$4 WHERE active = true AND userId = $1 AND movieId = $2`

	sqlRentSelect = `SELECT r.active, r.createAt, r.closeAt, m.id, m.name, m.year,
						array_to_string(movies.array_accum(g.name), ', ') AS genre, m.description `
	sqlRentCount = `SELECT COUNT(DISTINCT (r.id)) `
	sqlRentFrom  = `FROM movies.rent r, movies.movie m, movies.movie_genre mg, movies.genre g
						WHERE r.movieId = m.id AND mg.movieId = m.id AND mg.genreId = g.id
						AND r.userId = $1 AND r.active = $2 `
	sqlRentGroup = `GROUP BY r.id, m.id `
	sqlRentLimit = `LIMIT $3 OFFSET $4`
)

func (r Rent) Take() error {
	var CreateAt time.Time
	err := db.QueryRow(sqlRentValid, r.UserID, r.MovieID, true).Scan(&CreateAt)
	if err != nil && err == sql.ErrNoRows {
		err := r.create()
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("This movie you've already rented %v", CreateAt.Format("02-01-2006 15:04"))
}

func (r Rent) create() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	err = tx.QueryRow(sqlRentInsert, r.UserID, r.MovieID, true, r.CreateAt).Scan(&r.ID)
	if err != nil {
		return r.rowIsExist(err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r Rent) rowIsExist(err error) error {
	if strings.HasSuffix(err.Error(), "\"rent_movie_id_fk\"") {
		return fmt.Errorf("An identifier of the movie %v does not exist", r.MovieID)
	}
	return err
}

func (r Rent) GetTotalCount(userID int, active bool) (tc int, err error) {
	err = db.QueryRow(sqlRentCount+sqlRentFrom, userID, active).Scan(&tc)
	return
}

func (r Rent) GetAll(p *pagination.Pages, userID int, active bool) ([]RentMovie, error) {
	sql := strings.Join([]string{sqlRentSelect, sqlRentFrom, sqlRentGroup, sqlRentLimit}, "")
	rows, err := db.Query(sql, userID, active, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]RentMovie, 0)

	var CloseAt pq.NullTime
	for rows.Next() {

		rm := RentMovie{}
		err := rows.Scan(&rm.Active, &rm.CreateAt, &CloseAt, &rm.MovieID, &rm.Name, &rm.Year, &rm.Genre, &rm.Description)
		if err != nil {
			return nil, err
		}

		if CloseAt.Valid {
			rm.CloseAt = CloseAt.Time
		}
		res = append(res, rm)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Rent) Completed() error {
	err := db.QueryRow(sqlRentValid, r.UserID, r.MovieID, true).Scan(&r.CreateAt)
	if err != nil && err == sql.ErrNoRows {
		return fmt.Errorf("An identifier of the movie %v is not leased", r.MovieID)
	}
	if err != nil {
		return err
	}

	res, err := db.Exec(sqlRentUpdate, r.UserID, r.MovieID, false, time.Now())
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
