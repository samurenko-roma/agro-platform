package productionunit

import "errors"

var (
	ErrInvalidHierarchy = errors.New("invalid hierarchy")
	ErrAlreadyHasParent = errors.New("unit already has parent")
	ErrInvalidCode      = errors.New("invalid code")
	ErrAlreadyArchived  = errors.New("production unit already archived")
	ErrParentNotFound   = errors.New("parent production unit not found")
	ErrParentArchived   = errors.New("parent production unit is archived")
)
