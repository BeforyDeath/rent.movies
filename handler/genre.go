package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BeforyDeath/rent.movies/pagination"
	"github.com/BeforyDeath/rent.movies/request"
	"github.com/BeforyDeath/rent.movies/storage"
	"github.com/BeforyDeath/rent.movies/validator"
)

type genre struct{}

func (g genre) GetAll(w http.ResponseWriter, r *http.Request) {
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

	genre := new(storage.Genre)

	totalCount, err := genre.GetTotalCount()
	if err != nil {
		res.Error = err.Error()
		return
	}

	rows, err := genre.GetAll(pages)
	if err != nil {
		res.Error = err.Error()
		return
	}

	data := struct {
		Rows       []storage.Genre
		TotalCount int
	}{Rows: rows, TotalCount: totalCount}

	res.Success = true
	res.Data = data
	return
}
