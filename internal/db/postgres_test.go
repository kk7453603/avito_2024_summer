package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_GetPsqlDsn(t *testing.T) {
	cfg := &Config{
		Host:     "localhost",
		Port:     "5432",
		Name:     "postgres",
		User:     "postgres",
		Password: "password",
	}

	expected := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	actual := getPsqlDsn(cfg)

	require.Equal(t, expected, actual, "DSN must be correctly generated")
}
