package api

import (
	"net/http"
	"strconv"

	"github.com/juliankoehn/wetter-service/owm"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Response is the actual reponse of our API
// this is what we have to cache
type Response struct {
	Current *owm.CurrentWeatherData `json:"current"`
}

func (a *API) getWeather(c echo.Context) error {
	latitudeStr := c.QueryParam("lat")
	longitudeStr := c.QueryParam("lon")

	if latitudeStr == "" || longitudeStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing latitude or longitude")
	}
	// convert to float
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to parse given latitude")
	}
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to parse given longitude")
	}

	// read from cache
	key := getCacheKey(longitude, latitude)
	value, found := a.cache.Get(key)
	if !found {
		// lets try to get the closest cached item
		if a.config.OpenWeather.Fallback {
			return a.fallbackCall(c, longitude, latitude)
		}
		// if the default fallback is not allowed
		// we have to use our distance api
		// and measure to the closest point available in our cache
		return a.fallbackToClosest(c, longitude, latitude)

	}

	return c.JSON(200, value)
}

func (a *API) fallbackToClosest(c echo.Context, longitude, latitude float64) error {
	// we have to get the available locations from our cache
	// this is why we have stored a string[] map with our stores geolocations
	// to keep easy access without touching our storage engine (which is probably slow: sqlite)
	locations := a.getCachedLocations()
	if len(locations) == 0 {
		return echo.NewHTTPError(404, "No Data for given Location")
	}
	userPost := &owm.Coordinates{
		Longitude: longitude,
		Latitude:  latitude,
	}
	closest := locations[0]
	closestDistance := Distance(closest, userPost)

	for _, value := range locations {
		currDistance := Distance(value, userPost)
		if currDistance < closestDistance {
			closestDistance = currDistance
			closest = value
		}
	}

	value, found := a.readFromCache(closest.Longitude, closest.Latitude)
	if !found {
		logrus.Error("Error reading from Cache, Serving 404 to User. Entry SHOULD exists")
		return echo.NewHTTPError(404, "No Data for given Location")
	}

	return c.JSON(200, value)
}

// fallbackCall is used when owm fallback is enabled
// it gets executed when we have no entry in our cache
// this allows us to get data directly from owm and skipping the cache
func (a *API) fallbackCall(c echo.Context, longitude, latitude float64) error {
	res, err := owm.NewOneCall(owm.Celsius, owm.DE, a.config.OpenWeather.APIKEY, &owm.Coordinates{
		Longitude: longitude,
		Latitude:  latitude,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if ok := a.cacheByCoords(res, res.Longitude, res.Latitude); !ok {
		logrus.Error("Could not write Fallback item to cache")
	}
	// try to put into the cache

	return c.JSON(200, res)
}
