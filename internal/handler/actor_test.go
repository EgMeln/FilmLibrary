package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

type mockActorService struct {
	CreateFunc           func(actor *model.Actor) error
	UpdateFunc           func(actorID uuid.UUID, updatedActor *model.Actor) error
	DeleteFunc           func(actorID uuid.UUID) error
	GetAllWithMoviesFunc func() ([]*model.ActorMovies, error)
}

func (mas *mockActorService) Create(actor *model.Actor) error {
	return mas.CreateFunc(actor)
}

func (mas *mockActorService) Update(actorID uuid.UUID, updatedActor *model.Actor) error {
	return mas.UpdateFunc(actorID, updatedActor)
}

func (mas *mockActorService) Delete(actorID uuid.UUID) error {
	return mas.DeleteFunc(actorID)
}

func (mas *mockActorService) GetAllWithMovies() ([]*model.ActorMovies, error) {
	return mas.GetAllWithMoviesFunc()
}

func TestActorHandler_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		actor              model.Actor
		createFunc         func(actor *model.Actor) error
		expectedStatusCode int
	}{
		{
			name: "Success",
			actor: model.Actor{
				Name:   "name",
				Gender: "gender",
			},
			createFunc: func(actor *model.Actor) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "ServiceError",
			actor: model.Actor{
				Name:   "name",
				Gender: "gender",
			},
			createFunc: func(actor *model.Actor) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actorService := &mockActorService{
				CreateFunc: tc.createFunc,
			}
			actorHandler := NewActorHandler(actorService)

			jsonData, err := json.Marshal(tc.actor)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/actors/create", bytes.NewReader(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			actorHandler.Create(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestActorHandler_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		actorID            uuid.UUID
		updatedActor       model.Actor
		updateFunc         func(actorID uuid.UUID, updatedActor *model.Actor) error
		expectedStatusCode int
	}{
		{
			name:    "Success",
			actorID: uuid.New(),
			updatedActor: model.Actor{
				Name:   "name",
				Gender: "gender",
			},
			updateFunc: func(actorID uuid.UUID, updatedActor *model.Actor) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "ServiceError",
			actorID: uuid.New(),
			updatedActor: model.Actor{
				Name:   "updatedName",
				Gender: "updatedGender",
			},
			updateFunc: func(actorID uuid.UUID, updatedActor *model.Actor) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actorService := &mockActorService{
				UpdateFunc: tc.updateFunc,
			}
			actorHandler := NewActorHandler(actorService)

			jsonData, err := json.Marshal(tc.updatedActor)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, "/actors/update?actor_id="+tc.actorID.String(), bytes.NewReader(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			actorHandler.Update(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestActorHandler_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		actorID            uuid.UUID
		deleteFunc         func(actorID uuid.UUID) error
		expectedStatusCode int
	}{
		{
			name:    "Success",
			actorID: uuid.New(),
			deleteFunc: func(actorID uuid.UUID) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "ServiceError",
			actorID: uuid.New(),
			deleteFunc: func(actorID uuid.UUID) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actorService := &mockActorService{
				DeleteFunc: tc.deleteFunc,
			}
			actorHandler := NewActorHandler(actorService)

			req, err := http.NewRequest(http.MethodDelete, "/actors/delete?actor_id="+tc.actorID.String(), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			actorHandler.Delete(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestActorHandler_GetAllWithMovies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		getAllWithMoviesFunc func() ([]*model.ActorMovies, error)
		expectedStatusCode   int
	}{
		{
			name: "Success",
			getAllWithMoviesFunc: func() ([]*model.ActorMovies, error) {
				actorMovies := []*model.ActorMovies{
					{
						ID: uuid.New(), Name: "Actor1", Movies: []*model.Movie{
							{
								Title: "Movie1",
							},
							{
								Title: "Movie2",
							},
						},
					},
					{
						ID: uuid.New(), Name: "Actor2", Movies: []*model.Movie{
							{
								Title: "Movie3",
							},
							{
								Title: "Movie4",
							},
						},
					},
				}

				return actorMovies, nil
			},

			expectedStatusCode: http.StatusOK,
		},
		{
			name: "ServiceError",
			getAllWithMoviesFunc: func() ([]*model.ActorMovies, error) {
				return nil, errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actorService := &mockActorService{
				GetAllWithMoviesFunc: tc.getAllWithMoviesFunc,
			}
			actorHandler := NewActorHandler(actorService)

			req, err := http.NewRequest(http.MethodGet, "/actors/getAllWithMovies", nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			actorHandler.GetAllWithMovies(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}
