package storage

import "github.com/BeforyDeath/rent.movies/pagination"

type Genre struct {
	Id   int
	Name string
}

const (
	SQL_GENRE_LIST  = `SELECT id, name FROM movies.genre ORDER BY name ASC LIMIT $1 OFFSET $2`
	SQL_GENRE_COUNT = `SELECT COUNT(id) as c FROM movies.genre`
)

func (g Genre) GetTotalCount() (tc int, err error) {
	err = db.QueryRow(SQL_GENRE_COUNT).Scan(&tc)
	return
}

func (g Genre) GetAll(p *pagination.Pages) ([]Genre, error) {

	rows, err := db.Query(SQL_GENRE_LIST, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Genre, 0)
	for rows.Next() {
		genre := Genre{}
		err := rows.Scan(&genre.Id, &genre.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, genre)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
