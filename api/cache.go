package api

import (
	"fmt"
	"math"
	"time"

	"github.com/juliankoehn/wetter-service/owm"
	"github.com/sirupsen/logrus"
)

var (
	cacheKeyLocations = "cacheKeyLocations"
)

func (a *API) cacheByCoords(value *owm.OneCallResponse, long, lat float64) bool {
	key := getCacheKey(long, lat)

	// get the current location object
	locs, ok := a.cache.Get(cacheKeyLocations)
	if ok {
		oldLocs, k := locs.([]*owm.Coordinates)
		if k {
			replaced := false
			for index, location := range oldLocs {
				if location.Latitude == lat && location.Longitude == long {
					oldLocs[index] = &owm.Coordinates{
						Latitude:  lat,
						Longitude: long,
					}
					replaced = true
					break
				}
			}
			if !replaced {
				oldLocs = append(oldLocs, &owm.Coordinates{
					Latitude:  lat,
					Longitude: long,
				})
			}
			// store into cache
			if kk := a.cache.Set(cacheKeyLocations, oldLocs, 1); !kk {
				logrus.Error("Could not populate cacheKeyLocations")
			}
		}
	} else {
		// empty lets feed the cache
		newLocs := make([]*owm.Coordinates, 1)
		newLocs[0] = &owm.Coordinates{
			Latitude:  lat,
			Longitude: long,
		}
		// store into cache
		if kk := a.cache.Set(cacheKeyLocations, newLocs, 1); !kk {
			logrus.Error("Could not populate cacheKeyLocations")
		}

	}
	return a.writeToCache(key, value)
}

func (a *API) writeToCache(key string, value interface{}) bool {
	return a.cache.SetWithTTL(key, value, 1, 24*time.Hour)
}

func (a *API) readFromCache(long, lat float64) (interface{}, bool) {
	key := getCacheKey(long, lat)
	return a.cache.Get(key)
}

func (a *API) getCachedLocations() []*owm.Coordinates {
	val, ok := a.cache.Get(cacheKeyLocations)
	if !ok {
		return nil
	}
	// convert
	locs, ok := val.([]*owm.Coordinates)
	if !ok {
		return nil
	}

	return locs
}

// getCacheKey returns the cache key of our item
func getCacheKey(long, lat float64) string {
	// round to nearest
	// keeping 2 decimals
	// to match the OWM API
	long = math.Round(long*100) / 100
	lat = math.Round(lat*100) / 100

	return fmt.Sprintf("%f:%f", long, lat)
}
