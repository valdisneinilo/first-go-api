package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/go-sql-driver/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valdisneinilo/first-go-api/configs"
	_ "github.com/valdisneinilo/first-go-api/docs"
	"github.com/valdisneinilo/first-go-api/internal/infra/database"
	"github.com/valdisneinilo/first-go-api/internal/infra/webserver/handlers"
)

// @title           First GO API
// @version         1.0
// @description     Primeira API desenvolvida e GO.

// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization

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
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("Jwt", configs.JwtTokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JwtExpiresIn))

	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)
	userDB := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDB)

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

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	http.ListenAndServe(":8080", r)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}
