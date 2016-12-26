package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/BeforyDeath/rent.movies/config"
	"github.com/BeforyDeath/rent.movies/handler"
	"github.com/BeforyDeath/rent.movies/middleware"
	"github.com/BeforyDeath/rent.movies/storage"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// todo разобраться со структурой пакета базы, сделать пинг базы и реконект
	// todo сделать логер
	// todo перекрыть пакет errors, собрать свои коды ошибок

	filename := flag.String("f", "", "Initialized databases from filename")
	flag.Parse()

	cfg, err := config.NewConfig("config.json")
	if err != nil {
		fmt.Printf("JSON invalid file config: %v", err)
		return
	}

	store, err := storage.Connect("postgres", cfg.Database)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer store.Close()

	if *filename != "" {
		err = store.GetMigration(*filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Initialized databases from %v", *filename)
		return
	}

	// todo переделать
	handler.SetConfig(cfg)

	api := new(handler.Api)
	router := httprouter.New()

	api.Alice.NewChain(middleware.JsonContentType, middleware.FetchJsonRequest)

	router.POST("/user", api.Alice.AddChain(api.User.Create))
	router.POST("/login", api.Alice.AddChain(api.User.Login))

	router.POST("/genre", api.Alice.AddChain(api.Genre.GetAll))
	router.POST("/movie", api.Alice.AddChain(api.Movie.GetAll))

	router.POST("/rent/take", api.Alice.AddChain(api.Rent.Take, api.User.Authorization))
	router.POST("/rent/completed", api.Alice.AddChain(api.Rent.Completed, api.User.Authorization))
	router.POST("/rent/leased", api.Alice.AddChain(api.Rent.Leased, api.User.Authorization))

	fmt.Printf("Server started %s ...", cfg.Api.Listen)
	fmt.Println(http.ListenAndServe(cfg.Api.Listen, router))

}