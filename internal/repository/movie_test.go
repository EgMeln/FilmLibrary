package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

func TestMovieManager_Create(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)
}

func TestMovieManager_GetByID(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	getMovie, err := movieRep.GetByID(Barbi.ID)
	require.NoError(t, err)
	require.Equal(t, Barbi, getMovie)
}

func TestMovieManager_Update(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	updatedMovie := &model.Movie{
		ID:          Barbi.ID,
		Title:       "Oppenheimer",
		Description: "Atomic bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}

	err = movieRep.Update(updatedMovie)
	require.NoError(t, err)

	getMovie, err := movieRep.GetByID(Barbi.ID)
	require.NoError(t, err)
	require.Equal(t, updatedMovie, getMovie)
}

func TestMovieManager_Delete(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	Oppenheimer := &model.Movie{
		ID:          uuid.New(),
		Title:       "Oppenheimer",
		Description: "Atomic bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Oppenheimer)
	require.NoError(t, err)

	err = movieRep.Delete(Barbi.ID)
	require.NoError(t, err)

	getMovie, err := movieRep.GetByID(Barbi.ID)
	require.Error(t, err)
	require.Empty(t, getMovie)
}

func TestMovieManager_GetByTitle(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	Oppenheimer := &model.Movie{
		ID:          uuid.New(),
		Title:       "Oppenheimer",
		Description: "Atomic bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Oppenheimer)
	require.NoError(t, err)

	movies, err := movieRep.GetByTitle()
	require.NoError(t, err)
	require.Equal(t, []*model.Movie{Barbi, Oppenheimer}, movies)
}

func TestMovieManager_GetByRatingDesc(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      9,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	Oppenheimer := &model.Movie{
		ID:          uuid.New(),
		Title:       "Oppenheimer",
		Description: "Atomic bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Oppenheimer)
	require.NoError(t, err)

	movies, err := movieRep.GetByRatingDesc()
	require.NoError(t, err)
	require.Equal(t, []*model.Movie{Oppenheimer, Barbi}, movies)
}
func TestMovieManager_GetByReleaseDate(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      9,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	Oppenheimer := &model.Movie{
		ID:          uuid.New(),
		Title:       "Oppenheimer",
		Description: "Atomic bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Oppenheimer)
	require.NoError(t, err)

	movies, err := movieRep.GetByReleaseDate()
	require.NoError(t, err)
	require.Equal(t, []*model.Movie{Oppenheimer, Barbi}, movies)
}

func TestMovieManager_GetByTitleFragment(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}

	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      9,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	Oppenheimer := &model.Movie{
		ID:          uuid.New(),
		Title:       "Rbirb",
		Description: "bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Oppenheimer)
	require.NoError(t, err)

	movies, err := movieRep.GetByTitleFragment("rb")
	require.NoError(t, err)
	require.Equal(t, []*model.Movie{Barbi, Oppenheimer}, movies)
}

func TestMovieManager_GetByActorNameFragment(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movie_actor CASCADE")
		require.NoError(t, err)
		_, err = db.Exec("TRUNCATE TABLE movies CASCADE")
		require.NoError(t, err)
	}()

	Ken := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}
	err := actorRep.Create(Ken)
	require.NoError(t, err)

	Deadpool := &model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Reynolds ",
		Gender:    "Deadpool",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	}
	err = actorRep.Create(Deadpool)
	require.NoError(t, err)

	Barbi := &model.Movie{
		ID:          uuid.New(),
		Title:       "Barbi",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      9,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Barbi)
	require.NoError(t, err)

	Oppenheimer := &model.Movie{
		ID:          uuid.New(),
		Title:       "Oppenheimer",
		Description: "Atomic bomb",
		ReleaseDate: time.Date(2024, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Deadpool},
	}
	err = movieRep.Create(Oppenheimer)
	require.NoError(t, err)

	movies, err := movieRep.GetByActorNameFragment("Ryan")
	require.NoError(t, err)
	require.Equal(t, []*model.Movie{Barbi, Oppenheimer}, movies)
}
