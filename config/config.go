package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// OpenWeatherConfiguration keeps track of openweather informations
type OpenWeatherConfiguration struct {
	APIKEY           string `json:"api_key" required:"true" envconfig:"OPENWEATHER_API_KEY"`
	CityListEndpoint string `json:"city_list_endpoint" envconfig:"OPENWEATHER_CITYLIST_ENDPOINT"`
	Fallback         bool   `json:"fallback" envconfig:"OPENWEATHER_FALLBACK"`
}

// DBConfiguration connection params for our Database
type DBConfiguration struct {
	Driver    string `json:"driver" required:"true" envconfig:"DATABASE_DRIVER"`
	URL       string `json:"url" envconfig:"DATABASE_URL" required:"true"`
	Namespace string `json:"namespace" envconfig:"DATABASE_NAMESPACE"`
}

// WebConfiguration keeps information for our WebService
type WebConfiguration struct {
	UseTLS     bool   `json:"tls" envconfig:"WEB_USE_TLS"`
	ListenAddr string `json:"-"  envconfig:"WEB_LISTEN_ADDR"`
	BaseURL    string `json:"-"`
	Debug      bool   `json:"debug" envconfig:"WEB_DEBUG"`
}

// CacheConfiguration maps config for ristretto cache
type CacheConfiguration struct {
	NumCounters int64 `json:"num_counters" envconfig:"CACHE_NUM_COUNTERS"`
	MaxCosts    int64 `json:"max_costs" envconfig:"CACHE_MAX_COSTS"`
	BufferItems int64 `json:"buffer_items" envconfig:"CACHE_BUFFER_ITEMS"`
	Metrics     bool  `json:"metrics" envconfig:"CACHE_METRICS"`
}

// Configuration holds information about our current application-instance
// the config data SHOULD be loaded from environment variables
type Configuration struct {
	Logging     LoggingConfig            `json:"log"`
	DB          DBConfiguration          `json:"db"`
	OpenWeather OpenWeatherConfiguration `json:"open_weather"`
	Web         WebConfiguration         `json:"web"`
	Cache       CacheConfiguration       `json:"cache"`
}

// loadEnvironemnt checks if filename is present
// and loads config from filename OR default
func loadEnvironment(filename string) (err error) {
	if filename != "" {
		err = godotenv.Load(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exists, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

// LoadConfig loads the Configuration from env
func LoadConfig(filename string) (*Configuration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(Configuration)

	if err := envconfig.Process("weather", config); err != nil {
		return nil, err
	}

	if _, err := ConfigureLogging(&config.Logging); err != nil {
		return nil, err
	}
	config.applyDefaults()
	return config, nil
}

func (config *Configuration) applyDefaults() {
	// apply defaults to configration if empty
	if config.OpenWeather.CityListEndpoint == "" {
		config.OpenWeather.CityListEndpoint = "http://bulk.openweathermap.org/sample/city.list.json.gz"
	}

	if config.Cache.NumCounters == 0 {
		config.Cache.NumCounters = 1e7 // 10mb
	}
	if config.Cache.MaxCosts == 0 {
		config.Cache.MaxCosts = 1 << 30 // 1gb
	}

	if config.Cache.BufferItems == 0 {
		config.Cache.BufferItems = 64 // number of keys per Get buffer.
	}
}
