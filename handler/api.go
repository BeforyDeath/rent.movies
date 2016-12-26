package handler

import "github.com/BeforyDeath/rent.movies/config"

type API struct {
	User  user
	Movie movie
	Genre genre
	Rent  rent
	Alice alice
}

type Result struct {
	Success bool
	Data    interface{}
	Error   interface{}
}

var cfg *config.Сfg

func SetConfig(c *config.Сfg) {
	if c != nil {
		cfg = c
	}
}
