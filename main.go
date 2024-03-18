package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/EgMeln/filmLibraryPrivate/internal/config"
	"github.com/EgMeln/filmLibraryPrivate/internal/handler"
	"github.com/EgMeln/filmLibraryPrivate/internal/middleware"
	"github.com/EgMeln/filmLibraryPrivate/internal/repository"
	"github.com/EgMeln/filmLibraryPrivate/internal/service"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

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

	http.HandleFunc("/actors/create", middleware.AuthAdminMiddleware(actorHandler.Create))
	http.HandleFunc("/actors/update", middleware.AuthAdminMiddleware(actorHandler.Update))
	http.HandleFunc("/actors/delete", middleware.AuthAdminMiddleware(actorHandler.Delete))
	http.HandleFunc("/actors/getAllWithMovies", middleware.AuthUserMiddleware(actorHandler.GetAllWithMovies))

	http.HandleFunc("/movies/create", middleware.AuthAdminMiddleware(movieHandler.Create))
	http.HandleFunc("/movies/update", middleware.AuthAdminMiddleware(movieHandler.Update))
	http.HandleFunc("/movies/delete", middleware.AuthAdminMiddleware(movieHandler.Delete))
	http.HandleFunc("/movies/getAllWithSorting", middleware.AuthUserMiddleware(movieHandler.GetAllWithSorting))
	http.HandleFunc("/movies/getByTitleFragment", middleware.AuthUserMiddleware(movieHandler.GetByTitleFragment))
	http.HandleFunc("/movies/getByActorNameFragment", middleware.AuthUserMiddleware(movieHandler.GetByActorNameFragment))

	log.Printf("Server is running on %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, nil))
}
