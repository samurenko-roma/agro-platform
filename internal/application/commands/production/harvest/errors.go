package harvest

import "errors"

var (
	ErrHarvestNotFound      = errors.New("harvest batch not found")
	ErrGrowingCycleNotFound = errors.New("growing cycle not found")
)
