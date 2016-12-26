package config

import (
	"encoding/json"
	"os"
	"time"
)

type Сfg struct {
	Database string
	API      api
	Security security
}

type api struct {
	Listen    string
	PageLimit int64
}

type security struct {
	TokenSalt    string
	TokenExpired time.Duration
	PasswordSalt string
}

func NewConfig(path string) (*Сfg, error) {

	c := new(Сfg)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	if c.API.Listen == "" {
		c.API.Listen = ":8085"
	}

	return c, nil
}
