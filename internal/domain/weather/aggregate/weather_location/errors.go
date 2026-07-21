package weatherlocation

import "errors"

var (
	ErrLocationNotFound = errors.New("weather location not found")
	ErrLocationArchived = errors.New("weather location is archived")
)
