package api

import "github.com/labstack/echo/v4"

func (a *API) applyRoutes() {
	a.applyV1Routes()
}

// applyV1Routes applies the V1 route endpoint
// to our API
func (a *API) applyV1Routes() {
	v1 := a.e.Group("v1")
	v1.GET("*", a.getWeather)

	if a.config.Cache.Metrics {
		v1.GET("/metrics", a.cacheMetrics)
	}
}

func (a *API) cacheMetrics(c echo.Context) error {
	metrics := a.cache.Metrics

	return c.JSON(200, map[string]interface{}{
		"hits":       metrics.Hits(),
		"misses":     metrics.Misses(),
		"keys_added": metrics.KeysAdded(),
	})
}
