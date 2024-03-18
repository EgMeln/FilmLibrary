package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/EgMeln/filmLibraryPrivate/internal/service"
)

// MovieHandler handles HTTP requests related to movies.
type MovieHandler struct {
	movieService service.MovieService
}

// NewMovieHandler creates a new MovieHandler instance.
func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

// Create handles the HTTP request to create a new movie.
// @Summary Create a new movie
// @Description Create a new movie with the provided details
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body model.Movie true "Movie object to be created"
// @Success 200 {string} string "Movie created successfully"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to create movie"
// @Router /movies/create [post]
func (mh *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling Create Movie request...")

	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.movieService.Create(&movie); err != nil {
		http.Error(w, "Failed to create movie", http.StatusInternalServerError)
		log.Printf("Failed to create movie: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("Create Movie request handled successfully.")
}

// Update handles the HTTP request to update an existing movie.
// @Summary Update a movie
// @Description Update an existing movie with the provided details
// @Tags movies
// @Accept json
// @Produce json
// @Param movie_id query string true "ID of the movie to be updated"
// @Param movie body model.Movie true "Updated movie object"
// @Success 200 {string} string "Movie updated successfully"
// @Failure 400 {string} string "Invalid movie ID or failed to decode request body"
// @Failure 500 {string} string "Failed to update movie"
// @Router /movies/update [put]
func (mh *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling Update Movie request...")

	movieIDStr := r.URL.Query().Get("movie_id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		log.Printf("Invalid movie ID: %s", movieIDStr)
		return
	}

	var updatedMovie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.movieService.Update(movieID, updatedMovie); err != nil {
		http.Error(w, "Failed to update movie", http.StatusInternalServerError)
		log.Printf("Failed to update movie: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("Update Movie request handled successfully.")
}

// Delete handles the HTTP request to delete an existing movie.
// @Summary Delete a movie
// @Description Delete an existing movie by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Param movie_id query string true "ID of the movie to be deleted"
// @Success 200 {string} string "Movie deleted successfully"
// @Failure 400 {string} string "Invalid movie ID"
// @Failure 500 {string} string "Failed to delete movie"
// @Router /movies/delete [delete]
func (mh *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling Delete Movie request...")

	movieIDStr := r.URL.Query().Get("movie_id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		log.Printf("Invalid movie ID: %s", movieIDStr)
		return
	}

	if err := mh.movieService.Delete(movieID); err != nil {
		http.Error(w, "Failed to delete movie", http.StatusInternalServerError)
		log.Printf("Failed to delete movie: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("Delete Movie request handled successfully.")
}

// GetAllWithSorting handles the HTTP request to retrieve all movies with sorting.
// @Summary Get all movies with sorting
// @Description Retrieve all movies with sorting based on the provided flag
// @Tags movies
// @Accept json
// @Produce json
// @Param flag query int true "Sorting flag"
// @Success 200 {string} string "Movies retrieved successfully"
// @Failure 400 {string} string "Invalid sorting flag"
// @Failure 500 {string} string "Failed to fetch movies with sorting"
// @Router /movies/getAllWithSorting [get]
func (mh *MovieHandler) GetAllWithSorting(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling GetAllWithSorting Movies request...")

	flagStr := r.URL.Query().Get("flag")
	flag, err := strconv.Atoi(flagStr)
	if err != nil {
		http.Error(w, "Invalid sorting flag", http.StatusBadRequest)
		log.Printf("Invalid sorting flag: %s", flagStr)
		return
	}

	movies, err := mh.movieService.GetAllWithSorting(flag)
	if err != nil {
		http.Error(w, "Failed to fetch movies with sorting", http.StatusInternalServerError)
		log.Printf("Failed to fetch movies with sorting: %v", err)
		return
	}

	jsonResponse, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
		log.Printf("Failed to encode movies: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	log.Println("GetAllWithSorting Movies request handled successfully.")
}

// GetByTitleFragment handles the HTTP request to retrieve movies by title fragment.
// @Summary Get movies by title fragment
// @Description Retrieve movies that match the provided title fragment
// @Tags movies
// @Accept json
// @Produce json
// @Param title_fragment query string true "Title fragment"
// @Success 200 {string} string "Movies retrieved successfully"
// @Failure 500 {string} string "Failed to fetch movies by title fragment"
// @Router /movies/getByTitleFragment [get]
func (mh *MovieHandler) GetByTitleFragment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling GetByTitleFragment Movie request...")

	titleFragment := r.URL.Query().Get("title_fragment")

	movies, err := mh.movieService.GetByTitleFragment(titleFragment)
	if err != nil {
		http.Error(w, "Failed to fetch movies by title fragment", http.StatusInternalServerError)
		log.Printf("Failed to fetch movies by title fragment: %v", err)
		return
	}

	jsonResponse, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
		log.Printf("Failed to encode movies: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	log.Printf("GetByTitleFragment Movie request handled successfully.")
}

// GetByActorNameFragment handles the HTTP request to retrieve movies by actor name fragment.
// @Summary Get movies by actor name fragment
// @Description Retrieve movies associated with actors whose name matches the provided fragment
// @Tags movies
// @Accept json
// @Produce json
// @Param actor_name_fragment query string true "Actor name fragment"
// @Success 200 {string} string "Movies retrieved successfully"
// @Failure 500 {string} string "Failed to fetch movies by actor name fragment"
// @Router /movies/getByActorNameFragment [get]
func (mh *MovieHandler) GetByActorNameFragment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling GetByActorNameFragment Movie request...")

	actorNameFragment := r.URL.Query().Get("actor_name_fragment")

	movies, err := mh.movieService.GetByActorNameFragment(actorNameFragment)
	if err != nil {
		http.Error(w, "Failed to fetch movies by actor name fragment", http.StatusInternalServerError)
		log.Printf("Failed to fetch movies by actor name fragment: %v", err)
		return
	}

	jsonResponse, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
		log.Printf("Failed to encode movies: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	log.Printf("GetByActorNameFragment Movie request handled successfully.")
}
