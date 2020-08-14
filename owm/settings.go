package owm

import "net/http"

// Settings holds the client settings
type Settings struct {
	client *http.Client
}

// NewSettings returns a new Setting pointer with default http client.
func NewSettings() *Settings {
	return &Settings{
		client: http.DefaultClient,
	}
}
