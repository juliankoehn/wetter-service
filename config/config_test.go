package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

var (
	LogFile = "log.txt"
)

func TestMain(m *testing.M) {
	defer os.Clearenv()
	os.Exit(m.Run())
}

func TestGlobal(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FILE", LogFile)
	os.Setenv("LOG_DISABLE_COLORS", "true")
	os.Setenv("LOG_QUOTE_EMPTY_FIELDS", "true")
	os.Setenv("DATABASE_DRIVER", "sqlite3")
	os.Setenv("DATABASE_URL", "database.db")

	conf, err := LoadConfig("")
	require.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, "debug", conf.Logging.Level)
	assert.Equal(t, LogFile, conf.Logging.File)
	assert.Equal(t, true, conf.Logging.DisableColors)
	assert.Equal(t, true, conf.Logging.QuoteEmptyFields)

	// ignoring error, file could not be present if test fails
	_ = os.Remove(LogFile)
}
