package productionunit

import "errors"

var (
	ErrProductionUnitNotFound = errors.New("production unit not found")
	ErrHasActiveChildren      = errors.New("cannot archive: unit has active (non-archived) children")
)
