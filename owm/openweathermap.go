package owm

import "errors"

var (
	baseURL = "http://api.openweathermap.org/data/2.5/weather?%s"
	oneCall = "http://api.openweathermap.org/data/2.5/onecall?exclude=minutely&%s"
)

// DataUnit represents the character chosen to represent the temperature notation
type DataUnit string

const (
	// Fahrenheit represents the DataUnit imperial
	Fahrenheit DataUnit = "imperial"
	// Celsius represents the DataUnit metric
	Celsius = "metric"
)

// LangCode holds all supported languages to be used
// actually OMW is supporting way more languages
// but we only need Germany
type LangCode string

const (
	// DE is the German Language Code for OMW
	DE LangCode = "DE"
)

func setKey(key string) (string, error) {
	if err := ValidAPIKey(key); err != nil {
		return "", err
	}
	return key, nil
}

// ValidAPIKey makes sure that the key given is a valid one
func ValidAPIKey(key string) error {
	if len(key) != 32 {
		return errors.New("invalid key")
	}
	return nil
}
