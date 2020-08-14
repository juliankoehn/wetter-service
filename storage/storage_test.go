package storage

import (
	"os"
	"testing"

	"github.com/juliankoehn/wetter-service/config"
	"github.com/stretchr/testify/require"
)

var (
	testDBName = "database.test.db"
	testDBURI  = "sqlite3://database.test.db"
)

func TestStorage(t *testing.T) {
	config := &config.Configuration{
		DB: config.DBConfiguration{
			Driver: "sqlite3",
			URL:    testDBName,
		},
	}

	db, err := Connect(config)
	require.NoError(t, err)
	require.NotNil(t, db)

	// ignoring errors, if any
	os.Remove(testDBName)
}

func TestStorageWithURI(t *testing.T) {
	config := &config.Configuration{
		DB: config.DBConfiguration{
			URL: testDBURI,
		},
	}

	db, err := Connect(config)
	require.NoError(t, err)
	require.NotNil(t, db)

	// ignoring errors, if any
	os.Remove(testDBName)
}
