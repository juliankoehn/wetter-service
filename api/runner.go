package api

import (
	"time"

	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/owm"
	"github.com/sirupsen/logrus"
)

// this runs every x tick
// refreshes the data in cache of our
// weather data

const runnerDefaultTickTime = 5 * time.Minute

func (a *API) runner() {
	logrus.Info("Runner Cycle has started...")
	cities, err := model.FindEnabledCities(a.db)
	if err != nil {
		logrus.Errorf("error receiving enabled cities, skipping cycle: %+v", err)
		return
	}
	errors := 0

	logrus.Infof("Fetching Information for %d cities.", len(cities))
	for _, city := range cities {
		res, err := owm.NewOneCall(
			owm.Celsius,
			owm.DE,
			a.config.OpenWeather.APIKEY,
			&owm.Coordinates{
				Longitude: city.Coords.Longitude,
				Latitude:  city.Coords.Latitude,
			})
		if err != nil {
			errors++
			logrus.Errorf("Error fetching Info for %s - Error: %+v", city.Name, err)
			continue
		}

		//key := getCacheKey(city.Coords.Longitude, city.Coords.Latitude)
		// storing data for a maximum of 24 hours
		// will get overriden by cacheKey if we have new data
		if ok := a.cacheByCoords(res, city.Coords.Longitude, city.Coords.Latitude); !ok {
			errors++
			logrus.Errorf("Error caching item %s", city.Name)
		}
		//if ok := a.cache.SetWithTTL(key, res, 1, 24*time.Hour); !ok {
		//	errors++
		//	logrus.Error("Error caching item %s", city.Name)
		//}
	}
	logrus.Infof("Finished current Cycle, %d errors in this cycle", errors)
}
