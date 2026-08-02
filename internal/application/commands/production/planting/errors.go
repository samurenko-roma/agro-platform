package planting

import "errors"

var (
	ErrPlantingNotFound      = errors.New("planting not found")
	ErrInvalidProductionUnit = errors.New("invalid production unit")
	ErrCycleAlreadyExists    = errors.New("active cycle already exists")
	ErrGrowingCycleNotFound  = errors.New("growing cycle not found")
)
