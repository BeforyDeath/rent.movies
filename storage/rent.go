package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/BeforyDeath/rent.movies/pagination"
	"github.com/lib/pq"
)

type Rent struct {
	ID       int
	UserID   int
	MovieID  int64 `validate:"required"`
	Active   bool
	CreateAt time.Time
	CloseAt  time.Time
}

type RentList struct {
	Active   bool
	CreateAt time.Time
	CloseAt  time.Time
	MovieID  int

	Name        string
	Year        int
	Genre       string
	Description string
}

const (
	sqlRentValid  = `SELECT createAt FROM movies.rent WHERE userId = $1 AND movieId = $2 AND active = $3`
	sqlRentInsert = `INSERT INTO movies.rent (userId, movieId, active, createAt) VALUES ($1, $2, $3, $4) RETURNING id`
	sqlRentUpdate = `UPDATE movies.rent SET active=$3, closeAt=$4 WHERE active = true AND userId = $1 AND movieId = $2`

	sqlRentSelect = `SELECT r.active, r.createAt, r.closeAt, m.id as movieId, m.name, m.year,
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
	if err != nil {
		if err == sql.ErrNoRows {

			tx, err := db.Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			err = tx.QueryRow(sqlRentInsert, r.UserID, r.MovieID, true, r.CreateAt).Scan(&r.ID)
			if err != nil {
				if strings.HasSuffix(err.Error(), "\"rent_movie_id_fk\"") {
					return fmt.Errorf("An identifier of the movie %v does not exist", r.MovieID)
				}
				return err
			}

			err = tx.Commit()
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}
	return fmt.Errorf("This movie you've already rented %v", CreateAt.Format("02-01-2006 15:04"))
}

func (r Rent) GetTotalCount(userID int, active bool) (tc int, err error) {
	err = db.QueryRow(sqlRentCount+sqlRentFrom, userID, active).Scan(&tc)
	return
}

func (r Rent) GetAll(p *pagination.Pages, userID int, active bool) ([]RentList, error) {
	rows, err := db.Query(sqlRentSelect+sqlRentFrom+sqlRentGroup+sqlRentLimit, userID, active, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]RentList, 0)
	for rows.Next() {

		var CloseAt pq.NullTime

		rl := RentList{}
		err := rows.Scan(&rl.Active, &rl.CreateAt, &CloseAt, &rl.MovieID, &rl.Name, &rl.Year, &rl.Genre, &rl.Description)
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
	err := db.QueryRow(sqlRentValid, r.UserID, r.MovieID, true).Scan(&r.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("An identifier of the movie %v is not leased", r.MovieID)
		}
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
