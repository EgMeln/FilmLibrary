// Package repository - implements repository layer to interact with stored data in database
package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

// ActorManager represents an interface for managing actors in the system.
type ActorManager interface {
	Create(actor *model.Actor) error
	GetByID(actorID uuid.UUID) (*model.Actor, error)
	Update(actorID uuid.UUID, actor *model.Actor) error
	Delete(actorID uuid.UUID) error
	GetAllWithMovies() ([]*model.ActorMovies, error)
}

// NewActorStorage returns new repository instance for actors
func NewActorManager(db *sql.DB) ActorManager {
	return &actorManager{
		db: db,
	}
}

// actorManager implements crud method for actor
type actorManager struct {
	db *sql.DB
}

// Create inserts a new actor record into the database.
func (am *actorManager) Create(actor *model.Actor) error {
	query := `
		INSERT INTO actors (id, name, gender, birth_date) VALUES ($1, $2, $3, $4)`

	_, err := am.db.Query(query, actor.ID, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return err
	}
	return nil
}

// GetByID retrieves actor information from the database based on the provided actor ID.
func (am *actorManager) GetByID(actorID uuid.UUID) (*model.Actor, error) {
	query := `
		SELECT id, name, gender, birth_date AT TIME ZONE 'UTC' AS birth_date_utc 
		FROM actors 
		WHERE id = $1`

	var actor model.Actor

	err := am.db.QueryRow(query, actorID).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

// Update updates the information of an actor in the database.
func (am *actorManager) Update(actorID uuid.UUID, actor *model.Actor) error {
	query := `
		UPDATE actors SET name = COALESCE($2,name), gender = COALESCE($3,gender), birth_date = COALESCE($4,birth_date) 
		WHERE id = $1`

	_, err := am.db.Exec(query, actorID, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes actor information from the database based on the provided actor ID.
func (am *actorManager) Delete(actorID uuid.UUID) error {
	query := `DELETE FROM actors WHERE id = $1`

	_, err := am.db.Exec(query, actorID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllWithMovies retrieves a list of actors from the database along with information about the movies they starred in.
func (am *actorManager) GetAllWithMovies() ([]*model.ActorMovies, error) {
	query := `
	SELECT a.id AS actor_id, a.name AS actor_name, a.gender AS actor_gender, a.birth_date AT TIME ZONE 'UTC' AS actor_birth_date,
		   m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, m.release_date AT TIME ZONE 'UTC' AS movie_release_date, 
		   m.rating AS movie_rating
	FROM actors a
	JOIN movie_actor ma ON a.id = ma.actor_id
	JOIN movies m ON ma.movie_id = m.id`

	rows, err := am.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []*model.ActorMovies
	var actorMap = make(map[uuid.UUID]*model.ActorMovies)

	for rows.Next() {
		var nextActor model.ActorMovies
		var nextMovie model.Movie

		if err := rows.Scan(&nextActor.ID, &nextActor.Name, &nextActor.Gender, &nextActor.BirthDate,
			&nextMovie.ID, &nextMovie.Title, &nextMovie.Description, &nextMovie.ReleaseDate, &nextMovie.Rating); err != nil {
			return nil, err
		}

		actor, ok := actorMap[nextActor.ID]
		if !ok {
			actor = &nextActor
			actorMap[actor.ID] = actor
			actors = append(actors, actor)
		}

		actor.Movies = append(actor.Movies, &nextMovie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}