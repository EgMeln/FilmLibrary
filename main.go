package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/EgMeln/filmLibraryPrivate/internal/config"
	"github.com/EgMeln/filmLibraryPrivate/internal/handler"
	"github.com/EgMeln/filmLibraryPrivate/internal/middleware"
	"github.com/EgMeln/filmLibraryPrivate/internal/repository"
	"github.com/EgMeln/filmLibraryPrivate/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Unable to parse env: %v", err)
	}
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	actorManager := repository.NewActorManager(db)
	movieManager := repository.NewMovieManager(db)
	userManager := repository.NewUserManager(db)

	actorService := service.NewActorService(actorManager)
	movieService := service.NewMovieService(movieManager)
	userService := service.NewUserService(userManager)

	actorHandler := handler.NewActorHandler(actorService)
	movieHandler := handler.NewMovieHandler(movieService)
	userHandler := handler.NewUserHandler(userService)

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)

	http.HandleFunc("/actors/create", middleware.AuthMiddleware(actorHandler.Create))
	http.HandleFunc("/actors/update", middleware.AuthMiddleware(actorHandler.Update))
	http.HandleFunc("/actors/delete", middleware.AuthMiddleware(actorHandler.Delete))
	http.HandleFunc("/actors/getAllWithMovies", actorHandler.GetAllWithMovies)

	http.HandleFunc("/movies/create", middleware.AuthMiddleware(movieHandler.Create))
	http.HandleFunc("/movies/update", middleware.AuthMiddleware(movieHandler.Update))
	http.HandleFunc("/movies/delete", middleware.AuthMiddleware(movieHandler.Delete))
	http.HandleFunc("/movies/getAllWithSorting", movieHandler.GetAllWithSorting)
	http.HandleFunc("/movies/getByTitleFragment", movieHandler.GetByTitleFragment)
	http.HandleFunc("/movies/getByActorNameFragment", movieHandler.GetByActorNameFragment)

	log.Printf("Server is running on %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, nil))
}
