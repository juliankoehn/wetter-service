package owm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// utilizing the one-call-api
// makes it so fucking easy to retrieve all data and it is so inexpensive!

// Shared keeps shared units between responses
type Shared struct {
	Temp       float64   `json:"temp"`                // Temperature. Units - default: kelvin, metric: Celsius, imperial: Fahrenheit.
	FeelsLike  float64   `json:"feels_like"`          // Temperature. This temperature parameter accounts for the human perception of weather. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Pressure   float64   `json:"pressure"`            // Atmospheric pressure on the sea level, hPa
	Humidity   int       `json:"humidity"`            // Humidity, %
	DewPoint   float64   `json:"dew_point"`           // Atmospheric temperature (varying according to pressure and humidity) below which water droplets begin to condense and dew can form. Units – default: kelvin, metric: Celsius, imperial: Fahrenheit.
	Clouds     int       `json:"clouds"`              // Cloudiness, %
	Visibility int       `json:"visibility"`          // Average visibility, metres
	WindSpeed  float64   `json:"wind_speed"`          // Wind speed. Wind speed. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour
	WindGust   *float64  `json:"wind_gust,omitempty"` // (where available) Wind gust. Units – default: metre/sec, metric: metre/sec, imperial: miles/hour.
	WindDeg    float64   `json:"wind_deg"`            // Wind direction, degrees (meteorological)
	Rain       *Rain     `json:"rain,omitempty"`
	Snow       *Snow     `json:"snow,omitempty"`
	Weather    []Weather `json:"weather"`
}

// Current weather data API response
type Current struct {
	Time    int `json:"dt"`      //  Current time, Unix, UTC
	Sunrise int `json:"sunrise"` // Sunrise time, Unix, UTC
	Sunset  int `json:"sunset"`  // Sunset time, Unix, UTC
	Shared
	UV float64 `json:"uvi"` // UV Index
}

// Hourly keeps track of hourly weather changes
type Hourly struct {
	Time int `json:"dt"` //  Time of the forecasted data, Unix, UTC
	Shared
	Pop float64 `json:"pop"` // Probability of precipitation
}

// Daily keeps track of daily weather changes
type Daily struct {
	Time    int `json:"dt"`      //  Time of the forecasted data, Unix, UTC
	Sunrise int `json:"sunrise"` // Sunrise time, Unix, UTC
	Sunset  int `json:"sunset"`  // Sunset time, Unix, UTC
	Temp    struct {
		Morn  float64 `json:"morn"`  // Morning temperature.
		Day   float64 `json:"day"`   // day temperature.
		Eve   float64 `json:"eve"`   // Evening temperature.
		Night float64 `json:"night"` // Night temperature.
		Min   float64 `json:"min"`   // Min daily temperature.
		Max   float64 `json:"max"`   // Max daily temperature.
	} `json:"temp"`
	FeelsLike struct {
		Morn  float64 `json:"morn"`  // Morning temperature.
		Day   float64 `json:"day"`   // day temperature.
		Eve   float64 `json:"eve"`   // Evening temperature.
		Night float64 `json:"night"` // Night temperature.
	} `json:"feels_like"`
	Shared
	Rain float64 `json:"rain,omitempty"` // Precipitation volume, mm
	Snow float64 `json:"snow,omitempty"` // Snow volume, mm
	UV   float64 `json:"uvi,omitempty"`  // UV Index
}

// OneCallResponse is the response of the "oneCall" API from OWM
type OneCallResponse struct {
	Latitude       float64   `json:"lat"`
	Longitude      float64   `json:"lon"`
	Timezone       string    `json:"timezone"`
	TimezoneOffset int       `json:"timezone_offset"`
	Current        *Current  `json:"current"`
	Hourly         []*Hourly `json:"hourly"`
	Daily          []*Daily  `json:"daily"`
}

// NewOneCall returns a new OneCallResponse pointer with the supplied parameters
func NewOneCall(unit DataUnit, lang LangCode, key string, location *Coordinates) (*OneCallResponse, error) {
	if err := ValidAPIKey(key); err != nil {
		return nil, err
	}
	client := http.DefaultClient
	response, err := client.Get(fmt.Sprintf(fmt.Sprintf(oneCall, "appid=%s&lat=%f&lon=%f&units=%s&lang=%s"), key, location.Latitude, location.Longitude, unit, lang))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 200: is the only valid answer according to owm docs
	if response.StatusCode != http.StatusOK {
		return nil, handleResponseError(response)
	}
	res := &OneCallResponse{}
	if err = json.NewDecoder(response.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
