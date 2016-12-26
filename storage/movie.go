package storage

import (
	"database/sql"
	"strings"

	"github.com/BeforyDeath/rent.movies/pagination"
)

type Movie struct {
	ID          int
	Name        string
	Description string
	Year        int64  `validate:"neglect"`
	Genre       string `validate:"neglect"`
}

const (
	sqlMovieSelect = `SELECT m.id, m.name, m.year, array_to_string(movies.array_accum(g.name), ', ') as genre, m.description `
	sqlMovieCount  = `SELECT COUNT(DISTINCT (m.id)) `
	sqlMovieFrom   = `FROM movies.movie m, movies.movie_genre mg, movies.genre g WHERE mg.movieId = m.id AND mg.genreId = g.id `
	sqlMovieYear   = `AND m.year = $4 `
	sqlMovieGenre  = `AND m.id IN (SELECT m.id FROM movies.movie m, movies.movie_genre mg, movies.genre g
						WHERE mg.movieId = m.id AND mg.genreId = g.id AND g.name = $3 GROUP BY m.id) `
	sqlMovieGroup = `GROUP BY m.id `
	sqlMovieLimit = `ORDER BY m.id DESC LIMIT $1 OFFSET $2`
)

func (m Movie) ConstructSQL() map[string]string {
	built := make(map[string]string)

	var sqlBuilt, sqlBuiltCount string
	buildType := "limit"

	sqlBuilt += sqlMovieSelect
	sqlBuiltCount += sqlMovieCount

	sqlBuilt += sqlMovieFrom
	sqlBuiltCount += sqlMovieFrom

	if m.Year > 0 {
		sqlBuilt += sqlMovieYear
		sqlBuiltCount += sqlMovieYear
		buildType += "_year"
	}
	if m.Genre != "" {
		sqlBuilt += sqlMovieGenre
		sqlBuiltCount += sqlMovieGenre
		buildType += "_genre"
	}
	sqlBuilt += sqlMovieGroup
	sqlBuilt += sqlMovieLimit

	if buildType == "limit_year" {
		sqlBuilt = strings.Replace(sqlBuilt, "$4", "$3", 1)
	}
	if buildType == "limit_genre" {
		sqlBuiltCount = strings.Replace(sqlBuiltCount, "$3", "$1", 1)
	}
	sqlBuiltCount = strings.Replace(sqlBuiltCount, "$3", "$2", 1)
	sqlBuiltCount = strings.Replace(sqlBuiltCount, "$4", "$1", 1)

	built["type"] = buildType
	built["sql"] = sqlBuilt
	built["sql_count"] = sqlBuiltCount

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
		err := rows.Scan(&movie.ID, &movie.Name, &movie.Year, &movie.Genre, &movie.Description)
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
