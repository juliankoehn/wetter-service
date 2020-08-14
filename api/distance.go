package api

import (
	"math"

	"github.com/juliankoehn/wetter-service/owm"
)

const (
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

// as we dont want to store all geolocation data
// for each city in germany, we are using an algo
// to get the "nearest" location using Haversine

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Distance calculates the shortest path between two coordinates on the surface
// of the Earth.
func Distance(pos1, pos2 *owm.Coordinates) float64 {
	lat1 := degreesToRadians(pos1.Latitude)
	lon1 := degreesToRadians(pos1.Longitude)
	lat2 := degreesToRadians(pos2.Latitude)
	lon2 := degreesToRadians(pos2.Longitude)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km := c * earthRaidusKm

	return km
}
