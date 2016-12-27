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

	req := request.GetJSON(r)

	rent := new(storage.Rent)
	err := validator.GetRequest(req, rent)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := request.GetClaims(r)
	userID := int(claims["userID"].(float64))
	rent.UserID = userID

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

	req := request.GetJSON(r)

	rent := new(storage.Rent)
	err := validator.GetRequest(req, rent)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := request.GetClaims(r)
	userID := int(claims["userID"].(float64))
	rent.UserID = userID

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

	req := request.GetJSON(r)

	pages := new(pagination.Pages)
	err := validator.GetRequest(req, pages)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pages.Calculate(cfg.API.PageLimit)

	claims := request.GetClaims(r)
	userID := int(claims["userID"].(float64))

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
	totalCount, err := rent.GetTotalCount(userID, !show.History)
	if err != nil {
		res.Error = err.Error()
		return
	}

	rows, err := rent.GetAll(pages, userID, !show.History)
	if err != nil {
		res.Error = err.Error()
		return
	}

	data := struct {
		Rows       []storage.RentMovie
		TotalCount int
	}{Rows: rows, TotalCount: totalCount}

	res.Success = true
	res.Data = data
	return
}
