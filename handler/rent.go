package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BeforyDeath/rent.movies/pagination"
	"github.com/BeforyDeath/rent.movies/request"
	"github.com/BeforyDeath/rent.movies/storage"
	"github.com/BeforyDeath/rent.movies/validator"
)

type rent struct{}

func (re rent) Take(w http.ResponseWriter, r *http.Request) {
	res := Result{}
	defer func() { json.NewEncoder(w).Encode(res) }()

	req := request.GetJson(r)

	rent := new(storage.Rent)
	err := validator.GetRequest(req, rent)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := request.GetClaims(r)
	user_id := int(claims["user_id"].(float64))
	rent.User_id = user_id

	rent.CreateAt = time.Now()

	err = rent.Take()
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Success = true
	w.WriteHeader(http.StatusCreated)
	return
}

func (re rent) Completed(w http.ResponseWriter, r *http.Request) {
	res := Result{}
	defer func() { json.NewEncoder(w).Encode(res) }()

	req := request.GetJson(r)

	rent := new(storage.Rent)
	err := validator.GetRequest(req, rent)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := request.GetClaims(r)
	user_id := int(claims["user_id"].(float64))
	rent.User_id = user_id

	err = rent.Completed()
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Success = true
	return
}

func (re rent) Leased(w http.ResponseWriter, r *http.Request) {
	res := Result{}
	defer func() { json.NewEncoder(w).Encode(res) }()

	req := request.GetJson(r)

	pages := new(pagination.Pages)
	err := validator.GetRequest(req, pages)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pages.Calculate(cfg.Api.PageLimit)

	claims := request.GetClaims(r)
	user_id := int(claims["user_id"].(float64))

	var show = struct {
		History bool `validate:"neglect"`
	}{}
	err = validator.GetRequest(req, &show)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rent := new(storage.Rent)
	totalCount, err := rent.GetTotalCount(user_id, !show.History)
	if err != nil {
		res.Error = err.Error()
		return
	}

	rows, err := rent.GetAll(pages, user_id, !show.History)
	if err != nil {
		res.Error = err.Error()
		return
	}

	data := struct {
		Rows       []storage.RentList
		TotalCount int
	}{Rows: rows, TotalCount: totalCount}

	res.Success = true
	res.Data = data
	return
}
