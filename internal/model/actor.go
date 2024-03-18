// Package actors provides the Actor structure
package model

import (
	"time"

	"github.com/google/uuid"
)

// Actor represents information about an actor.
type Actor struct {
	ID        uuid.UUID // Unique identifier of the actor
	Name      string    // Name of the actor
	Gender    string    // Gender of the actor
	BirthDate time.Time // Birth date of the actor
}

// ActorMovies information about an actor and a list of films with his participation.
type ActorMovies struct {
	ID        uuid.UUID // Unique identifier of the actor
	Name      string    // Name of the actor
	Gender    string    // Gender of the actor
	BirthDate time.Time // Birth date of the actor
	Movies    []*Movie
}
