package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/google/uuid"
)

type mockMovieService struct {
	CreateFunc                 func(movie *model.Movie) error
	UpdateFunc                 func(movieID uuid.UUID, updatedMovie model.Movie) error
	DeleteFunc                 func(movieID uuid.UUID) error
	GetAllWithSortingFunc      func(flag int) ([]*model.Movie, error)
	GetByTitleFragmentFunc     func(titleFragment string) ([]*model.Movie, error)
	GetByActorNameFragmentFunc func(actorNameFragment string) ([]*model.Movie, error)
}

func (m *mockMovieService) Create(movie *model.Movie) error {
	return m.CreateFunc(movie)
}

func (m *mockMovieService) Update(movieID uuid.UUID, updatedMovie model.Movie) error {
	return m.UpdateFunc(movieID, updatedMovie)
}

func (m *mockMovieService) Delete(movieID uuid.UUID) error {
	return m.DeleteFunc(movieID)
}

func (m *mockMovieService) GetAllWithSorting(flag int) ([]*model.Movie, error) {
	return m.GetAllWithSortingFunc(flag)
}

func (m *mockMovieService) GetByTitleFragment(titleFragment string) ([]*model.Movie, error) {
	return m.GetByTitleFragmentFunc(titleFragment)
}

func (m *mockMovieService) GetByActorNameFragment(actorNameFragment string) ([]*model.Movie, error) {
	return m.GetByActorNameFragmentFunc(actorNameFragment)
}

func TestMovieHandler_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		movie              model.Movie
		createFunc         func(movie *model.Movie) error
		expectedStatusCode int
	}{
		{
			name: "Success",
			movie: model.Movie{
				Title: "Test Movie",
			},
			createFunc: func(movie *model.Movie) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "ServiceError",
			movie: model.Movie{
				Title: "Test Movie",
			},
			createFunc: func(movie *model.Movie) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMovieService{
				CreateFunc: tc.createFunc,
			}
			handler := NewMovieHandler(mockService)

			jsonData, err := json.Marshal(tc.movie)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/movies/create", bytes.NewReader(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			handler.Create(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestMovieHandler_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		movieID            uuid.UUID
		updatedMovie       model.Movie
		updateFunc         func(movieID uuid.UUID, updatedMovie model.Movie) error
		expectedStatusCode int
	}{
		{
			name:    "Success",
			movieID: uuid.New(),
			updatedMovie: model.Movie{
				Title: "Updated Movie",
			},
			updateFunc: func(movieID uuid.UUID, updatedMovie model.Movie) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "ServiceError",
			movieID: uuid.New(),
			updatedMovie: model.Movie{
				Title: "Updated Movie",
			},
			updateFunc: func(movieID uuid.UUID, updatedMovie model.Movie) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMovieService{
				UpdateFunc: tc.updateFunc,
			}
			handler := NewMovieHandler(mockService)

			jsonData, err := json.Marshal(tc.updatedMovie)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, "/movies/update?movie_id="+tc.movieID.String(), bytes.NewReader(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			handler.Update(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestMovieHandler_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		movieID            uuid.UUID
		deleteFunc         func(movieID uuid.UUID) error
		expectedStatusCode int
	}{
		{
			name:    "Success",
			movieID: uuid.New(),
			deleteFunc: func(movieID uuid.UUID) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "ServiceError",
			movieID: uuid.New(),
			deleteFunc: func(movieID uuid.UUID) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMovieService{
				DeleteFunc: tc.deleteFunc,
			}
			handler := NewMovieHandler(mockService)

			req, err := http.NewRequest(http.MethodDelete, "/movies/delete?movie_id="+tc.movieID.String(), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			handler.Delete(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestMovieHandler_GetAllWithSorting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		flag                  int
		getAllWithSortingFunc func(flag int) ([]*model.Movie, error)
		expectedStatusCode    int
	}{
		{
			name: "Success",
			flag: 1,
			getAllWithSortingFunc: func(flag int) ([]*model.Movie, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "ServiceError",
			flag: 1,
			getAllWithSortingFunc: func(flag int) ([]*model.Movie, error) {
				return nil, errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMovieService{
				GetAllWithSortingFunc: tc.getAllWithSortingFunc,
			}
			handler := NewMovieHandler(mockService)

			req, err := http.NewRequest(http.MethodGet, "/movies/getAllWithSorting?flag="+strconv.Itoa(tc.flag), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			handler.GetAllWithSorting(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestMovieHandler_GetByTitleFragment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		titleFragment          string
		getByTitleFragmentFunc func(titleFragment string) ([]*model.Movie, error)
		expectedStatusCode     int
	}{
		{
			name:          "Success",
			titleFragment: "fragment",
			getByTitleFragmentFunc: func(titleFragment string) ([]*model.Movie, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:          "ServiceError",
			titleFragment: "fragment",
			getByTitleFragmentFunc: func(titleFragment string) ([]*model.Movie, error) {
				return nil, errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMovieService{
				GetByTitleFragmentFunc: tc.getByTitleFragmentFunc,
			}
			handler := NewMovieHandler(mockService)

			req, err := http.NewRequest(http.MethodGet, "/movies/getByTitleFragment?title_fragment="+tc.titleFragment, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			handler.GetByTitleFragment(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestMovieHandler_GetByActorNameFragment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                       string
		actorNameFragment          string
		getByActorNameFragmentFunc func(actorNameFragment string) ([]*model.Movie, error)
		expectedStatusCode         int
	}{
		{
			name:              "Success",
			actorNameFragment: "fragment",
			getByActorNameFragmentFunc: func(actorNameFragment string) ([]*model.Movie, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:              "ServiceError",
			actorNameFragment: "fragment",
			getByActorNameFragmentFunc: func(actorNameFragment string) ([]*model.Movie, error) {
				return nil, errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMovieService{
				GetByActorNameFragmentFunc: tc.getByActorNameFragmentFunc,
			}
			handler := NewMovieHandler(mockService)

			req, err := http.NewRequest(http.MethodGet, "/movies/getByActorNameFragment?actor_name_fragment="+tc.actorNameFragment, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			handler.GetByActorNameFragment(recorder, req)
			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}
