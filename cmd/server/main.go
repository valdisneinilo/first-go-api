package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valdisneinilo/first-go-api/configs"
	"github.com/valdisneinilo/first-go-api/internal/infra/database"
	"github.com/valdisneinilo/first-go-api/internal/infra/webserver/handlers"
)

func main() {
	configs, err := configs.LoadConfigs(".")

	if err != nil {
		panic(err)
	}

	adress := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", configs.DbUser, configs.DbPassword, configs.DbHost, configs.DbPort, configs.DbName)
	db, err := sql.Open(configs.DbDriver, adress)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)

	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)
	userDB := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDB, configs.JwtTokenAuth, configs.JwtExpiresIn)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.JwtTokenAuth))
		r.Use(jwtauth.Authenticator(configs.JwtTokenAuth))

		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Patch("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Get("/", productHandler.GetAllProducts)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/signin", userHandler.GetJWT)

	http.ListenAndServe(":8080", r)
}
