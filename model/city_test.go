package model

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

var (
	testDbName = "database.db"
)

type CityTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (c *CityTestSuite) TestCity() {
	city := &City{
		ID:      "2803460",
		Name:    "Märkischer Kreis",
		State:   "",
		Country: "DE",
	}
	city.Coords.Latitude = 51.263889
	city.Coords.Longitude = 7.74167

	err := c.db.Create(city).Error
	assert.Equal(c.T(), err, nil)

	foundCity, err := FindCitiesByName(c.db, "Märkischer Kreis")
	assert.Equal(c.T(), err, nil)
	assert.NotEqual(c.T(), 0, len(foundCity))

	byID, err := FindCityByID(c.db, "2803460")
	assert.Equal(c.T(), err, nil)
	assert.NotNil(c.T(), byID.ID)

	// just copy oneonone city -> cityEnabled
	enabled := CityEnabled(*city)
	err = EnableCity(c.db, enabled)
	assert.Equal(c.T(), err, nil)

	ebyID, err := FindEnabledCityByID(c.db, "2803460")
	assert.Equal(c.T(), err, nil)
	assert.NotNil(c.T(), ebyID.ID)

	ec, err := FindEnabledCities(c.db)
	assert.Equal(c.T(), err, nil)
	assert.Equal(c.T(), len(ec), 1)

	foundCity = append(foundCity, &City{
		ID:      "2803468",
		Name:    "Zyfflich",
		State:   "",
		Country: "DE",
	})

	err = UpdateOrCreateCities(c.db, foundCity)

	assert.Equal(c.T(), err, nil)
}

func TestUserTestSuite(t *testing.T) {
	config := &config.Configuration{
		DB: config.DBConfiguration{
			Driver: "sqlite3",
			URL:    testDbName,
		},
	}
	db, err := storage.Connect(config)
	require.NoError(t, err)

	defer db.Close()
	db.AutoMigrate(&City{}, &CityEnabled{})

	ts := &CityTestSuite{
		db: db,
	}
	suite.Run(t, ts)

	// clean up db file from disk after test
	os.Remove(testDbName)
}
