package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BeforyDeath/rent.movies/pagination"
	"github.com/BeforyDeath/rent.movies/request"
	"github.com/BeforyDeath/rent.movies/storage"
	"github.com/BeforyDeath/rent.movies/validator"
)

type movie struct{}

func (m movie) GetAll(w http.ResponseWriter, r *http.Request) {
	res := Result{}
	defer func() { json.NewEncoder(w).Encode(res) }()

	req := request.GetJSON(r)

	pages := new(pagination.Pages)
	err := validator.GetRequest(req, pages)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pages.Calculate(cfg.API.PageLimit)

	movie := new(storage.Movie)
	err = validator.GetRequest(req, movie)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	built := movie.ConstructSQL()

	totalCount, err := movie.GetTotalCount(built)
	if err != nil {
		res.Error = err.Error()
		return
	}

	rows, err := movie.GetAll(built, pages)
	if err != nil {
		res.Error = err.Error()
		return
	}

	data := struct {
		Rows       []storage.Movie
		TotalCount int
	}{Rows: rows, TotalCount: totalCount}

	res.Success = true
	res.Data = data
	return
}
