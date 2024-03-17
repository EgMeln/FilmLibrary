package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

// MovieManager represents an interface for managing movies in the system.
type MovieManager interface {
	Create(movie *model.Movie) error
	GetByID(movieID uuid.UUID) (*model.Movie, error)
	Update(movie *model.Movie) error
	Delete(movieID uuid.UUID) error
	GetByTitle() ([]*model.Movie, error)
	GetByRatingDesc() ([]*model.Movie, error)
	GetByReleaseDate() ([]*model.Movie, error)
	GetByTitleFragment(fragment string) ([]*model.Movie, error)
	GetByActorNameFragment(fragment string) ([]*model.Movie, error)
}

// NewMovieManager returns new repository instance for movies
func NewMovieManager(db *sql.DB) MovieManager {
	return &movieManager{
		db: db,
	}
}

type movieManager struct {
	db *sql.DB
}

// Create inserts a new movie record along with its associated actors into the database.
func (mm *movieManager) Create(movie *model.Movie) error {
	tx, err := mm.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	movieQuery := `
		INSERT INTO movies (id, title, description, release_date, rating) VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(movieQuery, movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating)
	if err != nil {
		return err
	}

	actorQuery := `
		INSERT INTO movie_actor (movie_id, actor_id) VALUES ($1, $2)`

	for _, actor := range movie.Actors {
		_, err := tx.Exec(actorQuery, movie.ID, actor.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves movie information from the database based on the provided movie ID.
func (mm *movieManager) GetByID(movieID uuid.UUID) (*model.Movie, error) {
	movieQuery := `
        SELECT id, title, description, 
			release_date AT TIME ZONE 'UTC' AS release_date_utc, rating
        FROM movies
        WHERE id = $1
    `
	row := mm.db.QueryRow(movieQuery, movieID)

	var movie model.Movie
	err := row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	actorQuery := `
        SELECT a.id, a.name, a.gender, a.birth_date AT TIME ZONE 'UTC' AS birth_date_utc
        FROM actors a
        INNER JOIN movie_actor ma ON a.id = ma.actor_id
        WHERE ma.movie_id = $1
    `
	rows, err := mm.db.Query(actorQuery, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor model.Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return nil, err
		}
		movie.Actors = append(movie.Actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &movie, nil
}

// Update updates the information of a movie in the database based on the provided movie ID.
func (mm *movieManager) Update(movie *model.Movie) error {
	tx, err := mm.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	updateQuery := `
		UPDATE movies 
		SET title = COALESCE($2, title), 
			description = COALESCE($3,description), 
			release_date = COALESCE($4,release_date), 
			rating = COALESCE($5,rating) 
		WHERE id = $1`
	_, err = tx.Exec(updateQuery, movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating)
	if err != nil {
		return err
	}

	deleteQuery := `
		DELETE FROM movie_actor WHERE movie_id = $1`

	_, err = tx.Exec(deleteQuery, movie.ID)
	if err != nil {
		return err
	}

	actorQuery := `
		INSERT INTO movie_actor (movie_id, actor_id) VALUES ($1, $2)`

	for _, actor := range movie.Actors {
		_, err := tx.Exec(actorQuery, movie.ID, actor.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes movie information from the database based on the provided movie ID.
func (mm *movieManager) Delete(movieID uuid.UUID) error {
	tx, err := mm.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	q1 := `
		DELETE FROM movie_actor WHERE movie_id = $1`

	_, err = tx.Exec(q1, movieID)
	if err != nil {
		return err
	}

	q2 := `
		DELETE FROM movies WHERE id = $1`

	_, err = tx.Exec(q2, movieID)
	if err != nil {
		return err
	}

	return nil
}

// GetByTitle retrieves a list of movies from the database sorted by title,
func (mm *movieManager) GetByTitle() ([]*model.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.release_date AT TIME ZONE 'UTC' AS release_date_utc, m.rating,
			a.id, a.name, a.gender, a.birth_date AT TIME ZONE 'UTC' AS birth_date_utc
		FROM movies m
		LEFT JOIN movie_actor ma ON m.id = ma.movie_id
		LEFT JOIN actors a ON ma.actor_id = a.id
		ORDER BY m.title
`
	return mm.getMoviesByQuery(query)
}

// GetByRatingDesc retrieves a list of movies from the database sorted by rating
func (mm *movieManager) GetByRatingDesc() ([]*model.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.release_date AT TIME ZONE 'UTC' AS release_date_utc, m.rating,
			a.id, a.name, a.gender, a.birth_date AT TIME ZONE 'UTC' AS birth_date_utc 
		FROM movies m
		LEFT JOIN movie_actor ma ON m.id = ma.movie_id
		LEFT JOIN actors a ON ma.actor_id = a.id
		ORDER BY rating DESC
`
	return mm.getMoviesByQuery(query)
}

// GetByReleaseDate retrieves a list of movies from the database sorted by release date,
func (mm *movieManager) GetByReleaseDate() ([]*model.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.release_date AT TIME ZONE 'UTC' AS release_date_utc, m.rating,
			a.id, a.name, a.gender, a.birth_date AT TIME ZONE 'UTC' AS birth_date_utc
		FROM movies m
		LEFT JOIN movie_actor ma ON m.id = ma.movie_id
		LEFT JOIN actors a ON ma.actor_id = a.id
		ORDER BY m.release_date DESC
`
	return mm.getMoviesByQuery(query)
}

func (mm *movieManager) getMoviesByQuery(query string) ([]*model.Movie, error) {
	rows, err := mm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movieIDs []uuid.UUID
	movieMap := make(map[uuid.UUID]*model.Movie)

	for rows.Next() {
		var movie model.Movie
		var actor model.Actor

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating,
			&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return nil, err
		}

		if _, ok := movieMap[movie.ID]; !ok {
			movie.Actors = make([]model.Actor, 0)
			movieMap[movie.ID] = &movie
			movieIDs = append(movieIDs, movie.ID)
		} else {
			movie = *movieMap[movie.ID]
		}

		movie.Actors = append(movie.Actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var movies []*model.Movie
	for _, id := range movieIDs {
		movies = append(movies, movieMap[id])
	}

	return movies, nil
}

// GetByTitleFragment retrieves a list of movies from the database filtered by title fragment.
func (mm *movieManager) GetByTitleFragment(fragment string) ([]*model.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.release_date AT TIME ZONE 'UTC' AS release_date_utc, m.rating,
			a.id, a.name, a.gender, a.birth_date AT TIME ZONE 'UTC' AS birth_date_utc
		FROM movies m
		LEFT JOIN movie_actor ma ON m.id = ma.movie_id
		LEFT JOIN actors a ON ma.actor_id = a.id
		WHERE m.title LIKE '%' || $1 || '%'`

	return mm.getMoviesByQueryFragment(query, fragment)
}

// GetByActorNameFragment retrieves a list of movies from the database filtered by actor name fragment.
func (mm *movieManager) GetByActorNameFragment(fragment string) ([]*model.Movie, error) {
	query := `SELECT m.id, m.title, m.description, m.release_date AT TIME ZONE 'UTC' AS release_date_utc, m.rating,
		a.id, a.name, a.gender, a.birth_date AT TIME ZONE 'UTC' AS birth_date_utc
		FROM movies m
		LEFT JOIN movie_actor ma ON m.id = ma.movie_id
		LEFT JOIN actors a ON ma.actor_id = a.id
		WHERE a.name LIKE '%' || $1 || '%'`
	return mm.getMoviesByQueryFragment(query, fragment)
}

func (mm *movieManager) getMoviesByQueryFragment(query string, fragment string) ([]*model.Movie, error) {
	rows, err := mm.db.Query(query, fragment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movieIDs []uuid.UUID
	movieMap := make(map[uuid.UUID]*model.Movie)

	for rows.Next() {
		var movie model.Movie
		var actor model.Actor

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating,
			&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return nil, err
		}

		if _, ok := movieMap[movie.ID]; !ok {
			movie.Actors = make([]model.Actor, 0)
			movieMap[movie.ID] = &movie
			movieIDs = append(movieIDs, movie.ID)
		} else {
			movie = *movieMap[movie.ID]
		}

		movie.Actors = append(movie.Actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var movies []*model.Movie
	for _, id := range movieIDs {
		movies = append(movies, movieMap[id])
	}

	return movies, nil
}
