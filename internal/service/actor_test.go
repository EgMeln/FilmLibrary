package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

type mockActorManager struct {
	CreateFunc           func(actor *model.Actor) error
	GetByIDFunc          func(actorID uuid.UUID) (*model.Actor, error)
	UpdateFunc           func(actorID uuid.UUID, actor *model.Actor) error
	DeleteFunc           func(actorID uuid.UUID) error
	GetAllWithMoviesFunc func() ([]*model.ActorMovies, error)
}

func (m *mockActorManager) Create(actor *model.Actor) error {
	return m.CreateFunc(actor)
}

func (m *mockActorManager) GetByID(actorID uuid.UUID) (*model.Actor, error) {
	return m.GetByIDFunc(actorID)
}

func (m *mockActorManager) Update(actorID uuid.UUID, actor *model.Actor) error {
	return m.UpdateFunc(actorID, actor)
}

func (m *mockActorManager) Delete(actorID uuid.UUID) error {
	return m.DeleteFunc(actorID)
}

func (m *mockActorManager) GetAllWithMovies() ([]*model.ActorMovies, error) {
	return m.GetAllWithMoviesFunc()
}

func TestActorService_Create(t *testing.T) {
	t.Parallel()

	var (
		actors = make(map[string]*model.Actor)
	)

	mockManager := &mockActorManager{
		CreateFunc: func(actor *model.Actor) error {
			if _, exists := actors[actor.Name]; exists {
				return errors.New("actor already exists")
			}
			actors[actor.Name] = actor
			return nil
		},
	}

	tests := []struct {
		name           string
		actor          *model.Actor
		expectedResult error
	}{
		{
			name: "Success",
			actor: &model.Actor{
				ID:   uuid.New(),
				Name: "Tom Hanks",
			},
			expectedResult: nil,
		},
		{
			name: "Error",
			actor: &model.Actor{
				Name: "Tom Hanks",
			},
			expectedResult: errors.New("actor already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorSvc := &actorService{
				actorManager: mockManager,
			}

			err := actorSvc.Create(tt.actor)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestActorService_Update(t *testing.T) {
	t.Parallel()

	var (
		actors = make(map[string]*model.Actor)
	)

	mockManager := &mockActorManager{
		GetByIDFunc: func(actorID uuid.UUID) (*model.Actor, error) {
			actor, exists := actors[actorID.String()]
			if !exists {
				return nil, errors.New("actor not found")
			}
			return actor, nil
		},
		UpdateFunc: func(actorID uuid.UUID, actor *model.Actor) error {
			if _, exists := actors[actor.Name]; exists && actor.ID != actorID {
				return errors.New("actor with this name already exists")
			}
			actors[actor.Name] = actor
			return nil
		},
	}

	successActorID := uuid.New()
	actors[successActorID.String()] = &model.Actor{
		ID:   successActorID,
		Name: "Tom Hanks",
	}

	tests := []struct {
		name           string
		actorID        uuid.UUID
		actor          *model.Actor
		expectedResult error
	}{
		{
			name:    "Success",
			actorID: successActorID,
			actor: &model.Actor{
				Name: "Tom Hardy",
			},
			expectedResult: nil,
		},
		{
			name:    "ErrorActorNotFound",
			actorID: uuid.New(),
			actor: &model.Actor{
				Name: "Tom Hanks",
			},
			expectedResult: errors.New("actor not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorSvc := &actorService{
				actorManager: mockManager,
			}

			err := actorSvc.Update(tt.actorID, tt.actor)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestActorService_Delete(t *testing.T) {
	t.Parallel()

	var (
		actors = make(map[uuid.UUID]*model.Actor)
	)

	successActorID := uuid.New()
	actors[successActorID] = &model.Actor{
		ID:   successActorID,
		Name: "Tom Hanks",
	}

	mockManager := &mockActorManager{
		DeleteFunc: func(actorID uuid.UUID) error {
			if _, exists := actors[actorID]; !exists {
				return errors.New("actor not found")
			}
			delete(actors, actorID)
			return nil
		},
	}

	tests := []struct {
		name           string
		actorID        uuid.UUID
		expectedResult error
	}{
		{
			name:           "Success",
			actorID:        successActorID,
			expectedResult: nil,
		},
		{
			name:           "ErrorActorNotFound",
			actorID:        uuid.New(),
			expectedResult: errors.New("actor not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorSvc := &actorService{
				actorManager: mockManager,
			}

			err := actorSvc.Delete(tt.actorID)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestActorService_GetAllWithMovies(t *testing.T) {
	t.Parallel()
	actors := []*model.ActorMovies{
		{

			ID:   uuid.New(),
			Name: "Tom Hanks",

			Movies: []*model.Movie{
				{
					Title: "Forrest Gump",
				},
				{
					Title: "Saving Private Ryan",
				},
			},
		},
		{

			ID:   uuid.New(),
			Name: "Johnny Depp",

			Movies: []*model.Movie{
				{
					Title: "Pirates of the Caribbean",
				},
				{
					Title: "Edward Scissorhands",
				},
			},
		},
	}

	mockManager := &mockActorManager{
		GetAllWithMoviesFunc: func() ([]*model.ActorMovies, error) {
			return actors, nil
		},
	}

	actorSvc := &actorService{
		actorManager: mockManager,
	}

	expectedActors := []*model.ActorMovies{

		{

			ID:   actors[0].ID,
			Name: "Tom Hanks",

			Movies: []*model.Movie{
				{
					Title: "Forrest Gump",
				},
				{
					Title: "Saving Private Ryan",
				},
			},
		},
		{

			ID:   actors[1].ID,
			Name: "Johnny Depp",

			Movies: []*model.Movie{
				{
					Title: "Pirates of the Caribbean",
				},
				{
					Title: "Edward Scissorhands",
				},
			},
		},
	}

	actors, err := actorSvc.GetAllWithMovies()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(actors) != len(expectedActors) {
		t.Errorf("Expected %d actors, got %d", len(expectedActors), len(actors))
	}

	for i, actor := range actors {
		if actor.ID != expectedActors[i].ID || actor.Name != expectedActors[i].Name {
			t.Errorf("Expected actor: %+v, got: %+v", expectedActors[i], actor)
		}
		if len(actor.Movies) != len(expectedActors[i].Movies) {
			t.Errorf("Expected %d movies for actor %s, got %d", len(expectedActors[i].Movies), actor.Name, len(actor.Movies))
		}
	}
}
