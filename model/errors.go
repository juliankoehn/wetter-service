package model

// IsNotFoundError returns whether an error represents a "not found" error.
func IsNotFoundError(err error) bool {
	switch err.(type) {
	case CityNotFoundError:
		return true
	}
	return false
}

// CityNotFoundError represents when a user is not found.
type CityNotFoundError struct{}

func (e CityNotFoundError) Error() string {
	return "City not found"
}
