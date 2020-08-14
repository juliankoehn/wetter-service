package cities

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testCitiesEndpoint = "http://bulk.openweathermap.org/sample/city.list.json.gz"
)

func TestFileDownload(t *testing.T) {
	err := downloadUpdateFile(tempFileName, testCitiesEndpoint)
	require.NoError(t, err)
	// if no err a file named city.list.json.gz should exist in /cities
	ok := fileExists(tempFileName)
	require.True(t, ok)

	// remove file from storage
	_ = os.Remove(tempFileName)
}
