package service

import (
	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/EgMeln/filmLibraryPrivate/internal/repository"
)

// Sorting options for movies.
const (
	SortingByTitle       = 1
	SortingByReleaseDate = 2
)

// MovieService represents a service for managing movies.
type MovieService interface {
	Create(movie *model.Movie) error
	Update(movieID uuid.UUID, movie model.Movie) error
	Delete(movieID uuid.UUID) error
	GetAllWithSorting(flag int) ([]*model.Movie, error)
	GetByTitleFragment(titleFragment string) ([]*model.Movie, error)
	GetByActorNameFragment(actorNameFragment string) ([]*model.Movie, error)
}

type movieService struct {
	movieManager repository.MovieManager
}

// NewMovieService creates a new instance of the MovieService.
func NewMovieService(movieManager repository.MovieManager) MovieService {
	return &movieService{
		movieManager: movieManager,
	}
}

// Create creates a new movie.
func (ms *movieService) Create(movie *model.Movie) error {
	return ms.movieManager.Create(movie)
}

// Update updates an existing movie.
func (ms *movieService) Update(movieID uuid.UUID, movie model.Movie) error {
	existingMovie, err := ms.movieManager.GetByID(movieID)
	if err != nil {
		return err
	}

	if movie.Title != "" {
		existingMovie.Title = movie.Title
	}
	if movie.Description != "" {
		existingMovie.Description = movie.Description
	}
	if !movie.ReleaseDate.IsZero() {
		existingMovie.ReleaseDate = movie.ReleaseDate
	}
	if movie.Rating != 0 {
		existingMovie.Rating = movie.Rating
	}
	if movie.Actors != nil {
		existingMovie.Actors = movie.Actors
	}

	return ms.movieManager.Update(existingMovie)
}

// Delete deletes a movie by its ID.
func (ms *movieService) Delete(movieID uuid.UUID) error {
	return ms.movieManager.Delete(movieID)
}

// GetAllWithSorting retrieves all movies sorted by the specified flag.
func (ms *movieService) GetAllWithSorting(flag int) ([]*model.Movie, error) {
	var movies []*model.Movie
	var err error

	switch flag {
	case SortingByTitle:
		movies, err = ms.movieManager.GetByTitle()
	case SortingByReleaseDate:
		movies, err = ms.movieManager.GetByReleaseDate()
	default:
		movies, err = ms.movieManager.GetByRatingDesc()
	}

	if err != nil {
		return nil, err
	}
	return movies, err
}

// GetByTitleFragment retrieves movies containing the specified title fragment.
func (ms *movieService) GetByTitleFragment(titleFragment string) ([]*model.Movie, error) {
	return ms.movieManager.GetByTitleFragment(titleFragment)
}

// GetByActorNameFragment retrieves movies containing actors with the specified name fragment.
func (ms *movieService) GetByActorNameFragment(actorNameFragment string) ([]*model.Movie, error) {
	return ms.movieManager.GetByActorNameFragment(actorNameFragment)
}
