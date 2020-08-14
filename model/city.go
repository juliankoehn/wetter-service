package model

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/juliankoehn/wetter-service/storage/namespace"
	"github.com/pkg/errors"
)

// City reflects a City of openweathermap
// longest german city name is currently: Schmedeswurtherwesterdeich 26 chars.
type City struct {
	ID      json.Number `json:"id" gorm:"primary_key;unique_index"` // citylist can use floats as ID WTF
	Name    string      `json:"name" gorm:"type:varchar(100)"`
	State   string      `json:"state" gorm:"type:varchar(100)"` // SH / BW blabla
	Country string      `json:"country" gorm:"type:varchar(3)"` // max ISO3
	Coords  struct {
		Latitude  float64 `json:"lat" gorm:"column:lat"`
		Longitude float64 `json:"lon" gorm:"column:lon"`
	} `json:"coord" gorm:"embedded"`
}

// CityEnabled holds enabled cities
// we are going to use a second table `cities_enabled` as `cities` can
// get quite huge. the cities_enabled table is used to store "only" cities
// we are going to update
// this allows us fast selects without creating complex statistic tables (pg_statistics)
// we are not referencing *city to allow OneOnOne struct conversion
type CityEnabled struct {
	ID      json.Number `json:"id" gorm:"primary_key;unique_index"` // citylist can use floats as ID WTF
	Name    string      `json:"name" gorm:"type:varchar(100)"`
	State   string      `json:"state" gorm:"type:varchar(100)"` // SH / BW blabla
	Country string      `json:"country" gorm:"type:varchar(3)"` // max ISO3
	Coords  struct {
		Latitude  float64 `json:"lat" gorm:"column:lat"`
		Longitude float64 `json:"lon" gorm:"column:lon"`
	} `json:"coord" gorm:"embedded"`
}

// TableName returns the TableName of our user
func (ce *CityEnabled) TableName() string {
	tableName := "cities_enabled"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// TableName returns the TableName of our user
func (c *City) TableName() string {
	tableName := "cities"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// findCity finds a city in database
func findCity(tx *gorm.DB, query string, args ...interface{}) (*City, error) {
	obj := &City{}
	if err := tx.Where(query, args...).First(obj).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, CityNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding city")
	}

	return obj, nil
}

// findEnabledCity finds a city in database
func findEnabledCity(tx *gorm.DB, query string, args ...interface{}) (*City, error) {
	obj := &City{}
	if err := tx.Where(query, args...).First(obj).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, CityNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding city")
	}

	return obj, nil
}

// findEnabledCities the same as findCity but as a map
func findEnabledCities(tx *gorm.DB, query string, args ...interface{}) ([]*CityEnabled, error) {
	obj := []*CityEnabled{}
	if err := tx.Where(query, args...).Find(&obj).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, CityNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding cities")
	}
	return obj, nil
}

// FindEnabledCities returns all enabled cities
func FindEnabledCities(tx *gorm.DB) ([]*CityEnabled, error) {
	return findEnabledCities(tx, "", "")
}

// findCities the same as findCity but as a map
func findCities(tx *gorm.DB, query string, args ...interface{}) ([]*City, error) {
	obj := []*City{}
	if err := tx.Where(query, args...).Find(&obj).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, CityNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding cities")
	}
	return obj, nil
}

// FindCitiesByName searches our storage by name
// dont use fmt to inject % tags, this is a sec vuln in databases
func FindCitiesByName(tx *gorm.DB, name string) ([]*City, error) {
	pattern := "%" + name + "%"
	if tx.Dialect().GetName() == "sqlite3" {
		return findCities(tx, "name LIKE ? COLLATE NOCASE", pattern)
	}
	return findCities(tx, "name ILIKE ?", pattern)
}

// FindCityByID returns a City by it's id, error or nil
func FindCityByID(tx *gorm.DB, id string) (*City, error) {
	return findCity(tx, "id = ?", id)
}

// FindEnabledCityByID returns a City by it's id, error or nil
func FindEnabledCityByID(tx *gorm.DB, id string) (*City, error) {
	return findEnabledCity(tx, "id = ?", id)
}

// UpdateOrCreateCities updates or creates the given cities in our database
func UpdateOrCreateCities(tx *gorm.DB, cities []*City) error {
	for index, city := range cities {
		var c City
		// we have to iterate over each city
		// check if the entry already exists
		// if so update the entry
		// if not create the given entry
		if err := tx.Where("id = ?", city.ID).First(&c).Error; err != nil {
			// currently there are 3 possible errors
			// a) database is not setup
			// b) cannot connect to database
			// c) record is not found
			if gorm.IsRecordNotFoundError(err) {
				tx.Create(city) // create our city
			} else {
				return err
			}
		} else {
			// we have no error, so we found a city and we where able to bind it to c
			// let's try to update our city
			if err := tx.Model(cities[index]).
				Where("id = ?", city.ID).
				Update("name", city.Name).
				Update("lat", city.Coords.Latitude).
				Update("lon", city.Coords.Longitude).
				Error; err != nil {
				return err
			}
		}
	}
	// no error happened, just returning nil
	return nil
}

// EnableCity creates an entry in cities_enabled table
func EnableCity(tx *gorm.DB, city CityEnabled) error {
	// we are using firstOrCreate to prevent unique errors
	if err := tx.FirstOrCreate(&city).Error; err != nil {
		return err
	}
	return nil
}
