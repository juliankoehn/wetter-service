package owm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CurrentWeatherData struct contains an aggregate view of the structs
// defined above for JSON to be unmarshaled into.
type CurrentWeatherData struct {
	GeoPos   Coordinates `json:"coord"`
	Sys      Sys         `json:"sys"`
	Base     string      `json:"base"`
	Weather  []Weather   `json:"weather"`
	Main     Main        `json:"main"`
	Wind     Wind        `json:"wind"`
	Clouds   Clouds      `json:"clouds"`
	Rain     Rain        `json:"rain"`
	Snow     Snow        `json:"snow"`
	Dt       int         `json:"dt"`
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Cod      int         `json:"cod"`
	Timezone int         `json:"timezone"`
	Unit     DataUnit
	Lang     LangCode
	Key      string
	*Settings
}

// NewCurrent returns a new CurrentWeatherData pointer with the supplied parameters
func NewCurrent(unit DataUnit, lang LangCode, key string) (*CurrentWeatherData, error) {
	// instead of https://github.com/briandowns/openweathermap/blob/master/current.go
	// we can skip the entire validation stuff, as we are using constants
	// and no map interfaces....

	c := &CurrentWeatherData{
		Settings: NewSettings(),
		Unit:     unit,
		Lang:     lang,
	}

	var err error
	c.Key, err = setKey(key)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// CurrentByCoordinates will provide the current weather with the
// provided location coordinates.
func (w *CurrentWeatherData) CurrentByCoordinates(location *Coordinates) (*CurrentWeatherData, error) {
	response, err := w.client.Get(fmt.Sprintf(fmt.Sprintf(baseURL, "appid=%s&lat=%f&lon=%f&units=%s&lang=%s"), w.Key, location.Latitude, location.Longitude, w.Unit, w.Lang))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 200: is the only valid answer according to owm docs
	if response.StatusCode != http.StatusOK {
		return nil, handleResponseError(response)
	}

	if err = json.NewDecoder(response.Body).Decode(&w); err != nil {
		return nil, err
	}

	return w, nil
}
