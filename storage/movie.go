package storage

import (
	"database/sql"
	"github.com/BeforyDeath/rent.movies/pagination"
	"strings"
)

type Movie struct {
	Id          int
	Name        string
	Description string
	Year        int64  `validate:"neglect"`
	Genre       string `validate:"neglect"`
}

const (
	SQL_MOVIE_SELECT       = `SELECT m.id, m.name, m.year, array_to_string(movies.array_accum(g.name), ', ') as genre, m.description `
	SQL_MOVIE_SELECT_COUNT = `SELECT COUNT(DISTINCT (m.id)) `
	SQL_MOVIE_FROM         = `FROM movies.movie m, movies.movie_genre mg, movies.genre g WHERE mg.movie_id = m.id AND mg.genre_id = g.id `
	SQL_MOVIE_YEAR         = `AND m.year = $4 `
	SQL_MOVIE_GENRE        = `AND m.id IN (SELECT m.id FROM movies.movie m, movies.movie_genre mg, movies.genre g
						WHERE mg.movie_id = m.id AND mg.genre_id = g.id AND g.name = $3 GROUP BY m.id) `
	SQL_MOVIE_GROUP = `GROUP BY m.id `
	SQL_MOVIE_LIMIT = `ORDER BY m.id DESC LIMIT $1 OFFSET $2`
)

func (m Movie) ConstructSQL() map[string]string {
	built := make(map[string]string)

	var SQL_BUILT, SQL_BUILT_COUNT string
	buildType := "limit"

	SQL_BUILT += SQL_MOVIE_SELECT
	SQL_BUILT_COUNT += SQL_MOVIE_SELECT_COUNT

	SQL_BUILT += SQL_MOVIE_FROM
	SQL_BUILT_COUNT += SQL_MOVIE_FROM

	if m.Year > 0 {
		SQL_BUILT += SQL_MOVIE_YEAR
		SQL_BUILT_COUNT += SQL_MOVIE_YEAR
		buildType += "_year"
	}
	if m.Genre != "" {
		SQL_BUILT += SQL_MOVIE_GENRE
		SQL_BUILT_COUNT += SQL_MOVIE_GENRE
		buildType += "_genre"
	}
	SQL_BUILT += SQL_MOVIE_GROUP
	SQL_BUILT += SQL_MOVIE_LIMIT

	if buildType == "limit_year" {
		SQL_BUILT = strings.Replace(SQL_BUILT, "$4", "$3", 1)
	}
	if buildType == "limit_genre" {
		SQL_BUILT_COUNT = strings.Replace(SQL_BUILT_COUNT, "$3", "$1", 1)
	}
	SQL_BUILT_COUNT = strings.Replace(SQL_BUILT_COUNT, "$3", "$2", 1)
	SQL_BUILT_COUNT = strings.Replace(SQL_BUILT_COUNT, "$4", "$1", 1)

	built["type"] = buildType
	built["sql"] = SQL_BUILT
	built["sql_count"] = SQL_BUILT_COUNT

	return built
}

func (m Movie) GetTotalCount(built map[string]string) (tc int, err error) {

	switch built["type"] {
	case "limit":
		err = db.QueryRow(built["sql_count"]).Scan(&tc)
	case "limit_year":
		err = db.QueryRow(built["sql_count"], m.Year).Scan(&tc)
	case "limit_genre":
		err = db.QueryRow(built["sql_count"], m.Genre).Scan(&tc)
	case "limit_year_genre":
		err = db.QueryRow(built["sql_count"], m.Year, m.Genre).Scan(&tc)
	}
	if err != nil {
		return 0, err
	}
	return tc, nil
}

func (m Movie) GetAll(built map[string]string, p *pagination.Pages) ([]Movie, error) {

	var rows *sql.Rows
	var err error

	switch built["type"] {
	case "limit":
		rows, err = db.Query(built["sql"], p.Limit, p.Offset)
	case "limit_year":
		rows, err = db.Query(built["sql"], p.Limit, p.Offset, m.Year)
	case "limit_genre":
		rows, err = db.Query(built["sql"], p.Limit, p.Offset, m.Genre)
	case "limit_year_genre":
		rows, err = db.Query(built["sql"], p.Limit, p.Offset, m.Genre, m.Year)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Movie, 0)
	for rows.Next() {
		movie := Movie{}
		err := rows.Scan(&movie.Id, &movie.Name, &movie.Year, &movie.Genre, &movie.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
