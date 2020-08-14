package owm

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// APIError returned on failed API calls.
type APIError struct {
	Message string `json:"message"`
	COD     int    `json:"cod"`
}

// internal
type apiError struct {
	Message string      `json:"message"`
	COD     json.Number `json:"cod"`
}

func (ae *APIError) Error() string {
	return ae.Message
}

// 401: unauthorized or API key is not allowed to access endpoint
// 400: is user input error
func handleResponseError(r *http.Response) error {
	logrus.Info("Caught error handling!")
	var aError *apiError
	if err := json.NewDecoder(r.Body).Decode(&aError); err != nil {
		return err
	}

	// cast to APIERROR
	errorCode, err := strconv.Atoi(aError.COD.String())
	if err != nil {
		return err
	}

	return &APIError{
		Message: aError.Message,
		COD:     errorCode,
	}
}
