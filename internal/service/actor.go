// Package service provides implementations for various business logic services.
package service

import (
	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/EgMeln/filmLibraryPrivate/internal/repository"
)

// ActorService represents a service for managing actors.
type ActorService interface {
	Create(actor *model.Actor) error
	Update(actorID uuid.UUID, actor *model.Actor) error
	Delete(actorID uuid.UUID) error
	GetAllWithMovies() ([]*model.ActorMovies, error)
}

type actorService struct {
	actorManager repository.ActorManager
}

// NewActorService creates a new instance of the ActorService.
func NewActorService(actorManager repository.ActorManager) ActorService {
	return &actorService{
		actorManager: actorManager,
	}
}

// Create creates a new actor.
func (as *actorService) Create(actor *model.Actor) error {
	actor.ID = uuid.New()

	return as.actorManager.Create(actor)
}

// Update updates an existing actor.
func (as *actorService) Update(actorID uuid.UUID, actor *model.Actor) error {
	existingActor, err := as.actorManager.GetByID(actorID)
	if err != nil {
		return err
	}

	if actor.Name != "" {
		existingActor.Name = actor.Name
	}
	if actor.Gender != "" {
		existingActor.Gender = actor.Gender
	}
	if !actor.BirthDate.IsZero() {
		existingActor.BirthDate = actor.BirthDate
	}

	return as.actorManager.Update(actorID, existingActor)
}

// Delete deletes an actor by its ID.
func (as *actorService) Delete(actorID uuid.UUID) error {
	return as.actorManager.Delete(actorID)
}

// GetAllWithMovies retrieves all actors along with their movies.
func (as *actorService) GetAllWithMovies() ([]*model.ActorMovies, error) {
	return as.actorManager.GetAllWithMovies()
}
