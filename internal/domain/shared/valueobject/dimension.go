package valueobject

import "errors"

var ErrInvalidDimension = errors.New("invalid dimension")

type Dimension struct {
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
	Height float64 `json:"height"`
}

func NewDimension(w, l, h float64) (Dimension, error) {
	if w < 0 || l < 0 || h < 0 {
		return Dimension{}, ErrInvalidDimension
	}
	return Dimension{Width: w, Length: l, Height: h}, nil
}
