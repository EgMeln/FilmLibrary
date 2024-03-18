package model

import (
	"time"

	"github.com/google/uuid"
)

// Movie represents information about a movie.
type Movie struct {
    ID          uuid.UUID // Unique identifier of the movie
    Title       string    // Title of the movie
    Description string    // Description of the movie
    ReleaseDate time.Time // Release date of the movie
    Rating      int       // Rating of the movie
    Actors      []Actor   // List of actors starring in the movie
}
