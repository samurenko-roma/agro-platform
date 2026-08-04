package spatial

import (
	"github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
)

type provider struct {
	units     repository.ProductionUnitRepository
	snapshots repository.LayoutSnapshotRepository
}

func NewProvider() repository.SpatialProvider {
	return &provider{
		units:     NewProductionUnitRepository(),
		snapshots: NewLayoutSnapshotRepository(),
	}
}

func (p *provider) Units() repository.ProductionUnitRepository     { return p.units }
func (p *provider) Snapshots() repository.LayoutSnapshotRepository { return p.snapshots }
func (p *provider) ProviderName() string                           { return "spatial_inmemory" }

var _ repository.SpatialProvider = (*provider)(nil)
