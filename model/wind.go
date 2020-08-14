package model

// Wind holds information about the Speed and the direction the
// wind is moving.
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}
