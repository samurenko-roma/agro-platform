package repository

import (
	"context"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"

	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

type ProductionUnitRepository interface {
	Save(ctx context.Context, unit *pu.ProductionUnit) error
	GetByID(ctx context.Context, id vo.ID, orgId vo.ID) (*pu.ProductionUnit, error)
	GetChildren(ctx context.Context, parentID vo.ID) ([]*pu.ProductionUnit, error)
	ListByOwner(ctx context.Context, ownerID vo.ID) ([]*pu.ProductionUnit, error)
	Exists(ctx context.Context, id vo.ID) (bool, error)

	GetNextSequence(ctx context.Context, orgID vo.ID, parentID *vo.ID, unitType pu.ProductionUnitType) (int, error)
}
