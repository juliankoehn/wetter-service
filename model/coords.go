package model

// LatLng holds langitude and latitude
// this object is used in Requests and Responses
type LatLng struct {
	// The latitude in degrees. It must be in the range [-90.0, +90.0].
	Latitude float64 `json:"latitude"`
	// The longitude in degrees. It must be in the range [-180.0, +180.0].
	Longitude float64 `json:"longitude"`
}

// GetLatitude returns the Latitude if LatLng is present in mem or 0
func (x *LatLng) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

// GetLongitude returns the longitude if LatLng is present in mem or 0
func (x *LatLng) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}
