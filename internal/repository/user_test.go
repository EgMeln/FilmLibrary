package repository

import (
	"testing"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserManager_Create(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE users CASCADE")
		require.NoError(t, err)
	}()

	password, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	require.NoError(t, err)

	err = userRep.Create(&model.User{
		ID:       uuid.New(),
		Username: "admin",
		Password: string(password),
		Role:     "Ken",
	})
	require.NoError(t, err)
}

func TestUserManager_IfExist(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE users CASCADE")
		require.NoError(t, err)
	}()
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	require.NoError(t, err)
	user := &model.User{
		ID:       uuid.New(),
		Username: "admin",
		Password: string(password),
		Role:     "Ken",
	}

	err = userRep.Create(user)
	require.NoError(t, err)

	ifExist, err := userRep.IfExist(user.Username)
	require.NoError(t, err)
	require.Equal(t, true, ifExist)
}

func TestUserManager_GetByName(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE users CASCADE")
		require.NoError(t, err)
	}()
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	require.NoError(t, err)
	user := &model.User{
		ID:  uuid.New(),
		Username: "admin",
		Password: string(password),
		Role:     "Ken",
	}
	user.ID = uuid.New()
	err = userRep.Create(user)
	require.NoError(t, err)

	getUser, err := userRep.GetByUsername(user.Username)
	require.NoError(t, err)
	require.Equal(t, user, getUser)
}
