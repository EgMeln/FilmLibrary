package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

func TestActorManager_Create(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
		require.NoError(t, err)
	}()

	err := actorRep.Create(&model.Actor{
		ID:        uuid.New(),
		Name:      "Ryan Gosling",
		Gender:    "Drive",
		BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
}

func TestActorManager_GetByID(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
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

	getActor, err := actorRep.GetByID(Ken.ID)

	require.NoError(t, err)
	require.Equal(t, Ken, getActor)
}

func TestActorManager_Update(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
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

	updatedActor := &model.Actor{
		ID:        Ken.ID,
		Name:      "Cillian Murphy",
		Gender:    "Shelby",
		BirthDate: time.Date(1976, 5, 25, 0, 0, 0, 0, time.UTC),
	}

	err = actorRep.Update(Ken.ID, updatedActor)
	require.NoError(t, err)

	getActor, err := actorRep.GetByID(Ken.ID)
	require.NoError(t, err)
	require.Equal(t, updatedActor, getActor)
}

func TestActorManager_Delete(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE actors CASCADE")
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

	err = actorRep.Delete(Ken.ID)
	require.NoError(t, err)

	getActor, err := actorRep.GetByID(Ken.ID)
	require.Error(t, err)
	require.Empty(t, getActor)
}

func TestActorManager_GetAllWithMovies(t *testing.T) {
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
	Drive := &model.Movie{
		ID:          uuid.New(),
		Title:       "Drive",
		Description: "Ryan Gosling",
		ReleaseDate: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		Actors:      []model.Actor{*Ken},
	}
	err = movieRep.Create(Drive)
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

	actorsMovies, err := actorRep.GetAllWithMovies()
	require.NoError(t, err)

	expectedActorMovies :=
		[]*model.ActorMovies{
			{
				ID:        Ken.ID,
				Name:      "Ryan Gosling",
				Gender:    "Drive",
				BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
				Movies: []*model.Movie{
					{
						ID:          Barbi.ID,
						Title:       Barbi.Title,
						Description: Barbi.Description,
						ReleaseDate: Barbi.ReleaseDate,
						Rating:      Barbi.Rating,
					},
					{
						ID:          Drive.ID,
						Title:       Drive.Title,
						Description: Drive.Description,
						ReleaseDate: Drive.ReleaseDate,
						Rating:      Drive.Rating,
					},
				},
			},
			{
				ID:        Deadpool.ID,
				Name:      "Ryan Reynolds ",
				Gender:    "Deadpool",
				BirthDate: time.Date(1980, 11, 12, 0, 0, 0, 0, time.UTC),
				Movies: []*model.Movie{
					{
						ID:          Oppenheimer.ID,
						Title:       Oppenheimer.Title,
						Description: Oppenheimer.Description,
						ReleaseDate: Oppenheimer.ReleaseDate,
						Rating:      Oppenheimer.Rating,
					},
				},
			}}

	require.Equal(t, expectedActorMovies, actorsMovies)
}
