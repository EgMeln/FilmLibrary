package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

type mockMovieManager struct {
	CreateFunc                 func(movie *model.Movie) error
	GetByIDFunc                func(movieID uuid.UUID) (*model.Movie, error)
	UpdateFunc                 func(movie *model.Movie) error
	DeleteFunc                 func(movieID uuid.UUID) error
	GetByTitleFunc             func() ([]*model.Movie, error)
	GetByReleaseDateFunc       func() ([]*model.Movie, error)
	GetByRatingDescFunc        func() ([]*model.Movie, error)
	GetByTitleFragmentFunc     func(titleFragment string) ([]*model.Movie, error)
	GetByActorNameFragmentFunc func(actorNameFragment string) ([]*model.Movie, error)
}

func (m *mockMovieManager) Create(movie *model.Movie) error {
	return m.CreateFunc(movie)
}

func (m *mockMovieManager) GetByID(movieID uuid.UUID) (*model.Movie, error) {
	return m.GetByIDFunc(movieID)
}

func (m *mockMovieManager) Update(movie *model.Movie) error {
	return m.UpdateFunc(movie)
}

func (m *mockMovieManager) Delete(movieID uuid.UUID) error {
	return m.DeleteFunc(movieID)
}

func (m *mockMovieManager) GetByTitle() ([]*model.Movie, error) {
	return m.GetByTitleFunc()
}

func (m *mockMovieManager) GetByReleaseDate() ([]*model.Movie, error) {
	return m.GetByReleaseDateFunc()
}

func (m *mockMovieManager) GetByRatingDesc() ([]*model.Movie, error) {
	return m.GetByRatingDescFunc()
}

func (m *mockMovieManager) GetByTitleFragment(titleFragment string) ([]*model.Movie, error) {
	return m.GetByTitleFragmentFunc(titleFragment)
}

func (m *mockMovieManager) GetByActorNameFragment(actorNameFragment string) ([]*model.Movie, error) {
	return m.GetByActorNameFragmentFunc(actorNameFragment)
}

func TestMovieService_Create(t *testing.T) {
	t.Parallel()

	mockManager := &mockMovieManager{
		CreateFunc: func(movie *model.Movie) error {
			return nil
		},
	}

	tests := []struct {
		name           string
		movie          *model.Movie
		expectedResult error
	}{
		{
			name: "Success",
			movie: &model.Movie{
				ID:          uuid.New(),
				Title:       "Barbi",
				Description: "Barbie and Ken are having the time of their lives in the colorful and seemingly perfect world of Barbie Land.",
				ReleaseDate: time.Date(2024, time.July, 16, 0, 0, 0, 0, time.UTC),
				Rating:      8,
			},
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMovieService(mockManager)
			
			err := ms.Create(tt.movie)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestMovieService_Update(t *testing.T) {
	t.Parallel()

	mockManager := &mockMovieManager{
		GetByIDFunc: func(movieID uuid.UUID) (*model.Movie, error) {
			if movieID == uuid.Nil {
				return nil, errors.New("movie not found")
			}
			return &model.Movie{
				ID:          uuid.New(),
				Title:       "Barbi",
				Description: "Barbie and Ken are having the time of their lives in the colorful and seemingly perfect world of Barbie Land.",
				ReleaseDate: time.Date(2024, time.July, 16, 0, 0, 0, 0, time.UTC),
				Rating:      8,
			}, nil
		},
		UpdateFunc: func(movie *model.Movie) error {
			return nil
		},
	}

	tests := []struct {
		name           string
		movieID        uuid.UUID
		movie          model.Movie
		expectedResult error
	}{
		{
			name:    "Success",
			movieID: uuid.New(),
			movie: model.Movie{
				ID:          uuid.New(),
				Title:       "Barbi",
				Description: "Barbie and Ken are having the time of their lives in the colorful and seemingly perfect world of Barbie Land.",
				ReleaseDate: time.Date(2024, time.July, 16, 0, 0, 0, 0, time.UTC),
				Rating:      8,
				Actors: []model.Actor{
					{
						Name: "Tom Hanks",
					},
				},
			},
			expectedResult: nil,
		},
		{
			name:    "MovieNotFound",
			movieID: uuid.Nil,
			movie: model.Movie{
				ID:          uuid.New(),
				Title:       "Barbi",
				Description: "Barbie and Ken are having the time of their lives in the colorful and seemingly perfect world of Barbie Land.",
				ReleaseDate: time.Date(2024, time.July, 16, 0, 0, 0, 0, time.UTC),
				Rating:      8,
			},
			expectedResult: errors.New("movie not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMovieService(mockManager)

			err := ms.Update(tt.movieID, tt.movie)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestMovieService_Delete(t *testing.T) {
	t.Parallel()

	mockManager := &mockMovieManager{
		DeleteFunc: func(movieID uuid.UUID) error {
			return nil
		},
	}

	tests := []struct {
		name           string
		movieID        uuid.UUID
		expectedResult error
	}{
		{
			name:           "Success",
			movieID:        uuid.New(),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMovieService(mockManager)

			err := ms.Delete(tt.movieID)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestMovieService_GetAllWithSorting(t *testing.T) {
	t.Parallel()

	mockManager := &mockMovieManager{
		GetByTitleFunc: func() ([]*model.Movie, error) {
			return []*model.Movie{
				{Title: "Movie1"},
				{Title: "Movie2"},
			}, nil
		},
		GetByReleaseDateFunc: func() ([]*model.Movie, error) {
			return []*model.Movie{
				{Title: "Movie1", ReleaseDate: time.Date(2024, time.July, 16, 0, 0, 0, 0, time.UTC)},
				{Title: "Movie2", ReleaseDate: time.Date(2023, time.July, 16, 0, 0, 0, 0, time.UTC)},
			}, nil
		},
		GetByRatingDescFunc: func() ([]*model.Movie, error) {
			return []*model.Movie{
				{Title: "Movie2", Rating: 8},
				{Title: "Movie1", Rating: 7},
			}, nil
		},
	}

	tests := []struct {
		name           string
		flag           int
		expectedResult []*model.Movie
	}{
		{
			name:           "SortingByTitle",
			flag:           SortingByTitle,
			expectedResult: []*model.Movie{{Title: "Movie1"}, {Title: "Movie2"}},
		},
		{
			name:           "SortingByReleaseDate",
			flag:           SortingByReleaseDate,
			expectedResult: []*model.Movie{{Title: "Movie1", ReleaseDate: time.Date(2024, time.July, 16, 0, 0, 0, 0, time.UTC)}, {Title: "Movie2", ReleaseDate: time.Date(2023, time.July, 16, 0, 0, 0, 0, time.UTC)}},
		},
		{
			name:           "SortingByRating",
			expectedResult: []*model.Movie{{Title: "Movie2", Rating: 8}, {Title: "Movie1", Rating: 7}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMovieService(mockManager)

			movies, err := ms.GetAllWithSorting(tt.flag)

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if len(movies) != len(tt.expectedResult) {
				t.Errorf("Expected %d movies, got: %d", len(tt.expectedResult), len(movies))
			}

			for i := range movies {
				if movies[i].Title != tt.expectedResult[i].Title {
					t.Errorf("Expected movie title: %s, got: %s", tt.expectedResult[i].Title, movies[i].Title)
				}
			}
		})
	}
}

func TestMovieService_GetByTitleFragment(t *testing.T) {
	t.Parallel()

	mockManager := &mockMovieManager{
		GetByTitleFragmentFunc: func(titleFragment string) ([]*model.Movie, error) {
			return []*model.Movie{
				{Title: "Movie1"},
				{Title: "Movie2"},
			}, nil
		},
	}

	tests := []struct {
		name           string
		titleFragment  string
		expectedResult []*model.Movie
	}{
		{
			name:           "Fragment1",
			titleFragment:  "frag",
			expectedResult: []*model.Movie{{Title: "Movie1"}, {Title: "Movie2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMovieService(mockManager)

			movies, err := ms.GetByTitleFragment(tt.titleFragment)

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if len(movies) != len(tt.expectedResult) {
				t.Errorf("Expected %d movies, got: %d", len(tt.expectedResult), len(movies))
			}

			for i := range movies {
				if movies[i].Title != tt.expectedResult[i].Title {
					t.Errorf("Expected movie title: %s, got: %s", tt.expectedResult[i].Title, movies[i].Title)
				}
			}
		})
	}
}
func TestMovieService_GetByActorNameFragment(t *testing.T) {
	t.Parallel()

	mockManager := &mockMovieManager{
		GetByActorNameFragmentFunc: func(actorNameFragment string) ([]*model.Movie, error) {
			return []*model.Movie{
				{Title: "Movie1"},
				{Title: "Movie2"},
			}, nil
		},
	}

	tests := []struct {
		name              string
		actorNameFragment string
		expectedResult    []*model.Movie
	}{
		{
			name:              "Fragment1",
			actorNameFragment: "frag",
			expectedResult:    []*model.Movie{{Title: "Movie1"}, {Title: "Movie2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMovieService(mockManager)

			movies, err := ms.GetByActorNameFragment(tt.actorNameFragment)

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if len(movies) != len(tt.expectedResult) {
				t.Errorf("Expected %d movies, got: %d", len(tt.expectedResult), len(movies))
			}

			for i := range movies {
				if movies[i].Title != tt.expectedResult[i].Title {
					t.Errorf("Expected movie title: %s, got: %s", tt.expectedResult[i].Title, movies[i].Title)
				}
			}
		})
	}
}
