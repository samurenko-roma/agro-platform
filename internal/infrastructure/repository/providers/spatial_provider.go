package providers

import (
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	inmemory "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/inmemory/spatial"
	pgsnapshot "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/spatial/layout_snapshot"
	postgres "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/spatial/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type spatialProvider struct {
	db        uow.DB
	inMemory  bool
	units     domain.ProductionUnitRepository
	snapshots domain.LayoutSnapshotRepository
}

func (p *spatialProvider) ProviderName() string {
	return "spatial"
}

func NewSpatialProvider(db uow.DB) repository.RepositoryProvider {
	return &spatialProvider{db: db, inMemory: false}
}

func (p *spatialProvider) Units() domain.ProductionUnitRepository {
	if p.units != nil {
		return p.units
	}
	if p.inMemory {
		p.units = inmemory.NewProductionUnitRepository()
	} else {
		p.units = postgres.New(p.db)
	}
	return p.units
}

func (p *spatialProvider) Snapshots() domain.LayoutSnapshotRepository {
	if p.snapshots != nil {
		return p.snapshots
	}
	if p.inMemory {
		p.snapshots = inmemory.NewLayoutSnapshotRepository()
	} else {
		p.snapshots = pgsnapshot.NewRepository(p.db)
	}
	return p.snapshots
}

var _ domain.SpatialProvider = (*spatialProvider)(nil)
