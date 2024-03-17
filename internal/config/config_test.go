package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	var environment = map[string]interface{}{
		"POSTGRES_URL": "postgresql://postgres:a!11111111@localhost:5432/profile",
		"SERVER_PORT":  "localhost:5432",
	}

	for k, v := range environment {
		require.Nil(t, os.Setenv(k, fmt.Sprintf("%v", v)))
	}

	cfg, err := NewConfig()
	require.Nil(t, err)

	require.Equal(t, environment["SERVER_PORT"], cfg.ServerPort)
	require.Equal(t, environment["POSTGRES_URL"], cfg.PostgresURL)
}
