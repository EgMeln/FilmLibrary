// Package handler provides HTTP handlers for managing in the film library.
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/EgMeln/filmLibraryPrivate/internal/service"
)

// ActorHandler handles HTTP requests related to actors.
type ActorHandler struct {
	actorService service.ActorService
}

// NewActorHandler creates a new ActorHandler instance.
func NewActorHandler(actorService service.ActorService) *ActorHandler {
	return &ActorHandler{
		actorService: actorService,
	}
}

// Create handles HTTP requests to create a new actor.
//	@Summary		Create a new actor
//	@Description	Create a new actor in the film library
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			actor	body		model.Actor	true	"Actor object to be created"
//	@Success		200		{string}	string		"OK"
//	@Failure		400		{string}	string		"Failed to decode request body"
//	@Failure		500		{string}	string		"Failed to create actor"
//	@Router			/actors/create [post]
func (ah *ActorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var actor model.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := ah.actorService.Create(&actor); err != nil {
		http.Error(w, "Failed to create actor", http.StatusInternalServerError)
		log.Printf("Failed to create actor: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Update handles HTTP requests to update an existing actor.
//	@Summary		Update an existing actor
//	@Description	Update an existing actor in the film library
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			actor_id	query		string		true	"ID of the actor to be updated"
//	@Param			actor		body		model.Actor	true	"Actor object with updated information"
//	@Success		200			{string}	string		"OK"
//	@Failure		400			{string}	string		"Invalid actor ID"
//	@Failure		400			{string}	string		"Failed to decode request body"
//	@Failure		500			{string}	string		"Failed to update actor"
//	@Router			/actors/update [put]
func (ah *ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("actor_id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		log.Printf("Invalid actor ID: %s", actorIDStr)
		return
	}

	var updatedActor model.Actor
	if err := json.NewDecoder(r.Body).Decode(&updatedActor); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := ah.actorService.Update(actorID, &updatedActor); err != nil {
		http.Error(w, "Failed to update actor", http.StatusInternalServerError)
		log.Printf("Failed to update actor: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete handles HTTP requests to delete an actor by ID.
//	@Summary		Delete an actor
//	@Description	Delete an actor from the film library by ID
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Param			actor_id	query		string	true	"ID of the actor to be deleted"
//	@Success		200			{string}	string	"OK"
//	@Failure		400			{string}	string	"Invalid actor ID"
//	@Failure		500			{string}	string	"Failed to delete actor"
//	@Router			/actors/delete [delete]
func (ah *ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("actor_id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		log.Printf("Invalid actor ID: %s", actorIDStr)
		return
	}

	if err := ah.actorService.Delete(actorID); err != nil {
		http.Error(w, "Failed to delete actor", http.StatusInternalServerError)
		log.Printf("Failed to delete actor: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetAllWithMovies handles HTTP requests to retrieve all actors with their associated movies.
//	@Summary		Retrieve all actors with associated movies
//	@Description	Retrieve all actors from the film library along with their associated movies
//	@Tags			actors
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]model.ActorMovies	"OK"
//	@Failure		500	{string}	string				"Failed to fetch actors with movies"
//	@Router			/actors/getAllWithMovies [get]
func (ah *ActorHandler) GetAllWithMovies(w http.ResponseWriter, r *http.Request) {
	actorMovies, err := ah.actorService.GetAllWithMovies()
	if err != nil {
		http.Error(w, "Failed to fetch actors with movies", http.StatusInternalServerError)
		log.Printf("Failed to fetch actors with movies: %v", err)
		return
	}

	jsonResponse, err := json.Marshal(actorMovies)
	if err != nil {
		http.Error(w, "Failed to encode actors with movies", http.StatusInternalServerError)
		log.Printf("Failed to encode actors with movies: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
