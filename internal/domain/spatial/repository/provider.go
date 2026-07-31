package repository

import "github.com/samurenkoroma/agro-platform/internal/shared/repository"

type SpatialProvider interface {
	repository.RepositoryProvider
	Units() ProductionUnitRepository
	Snapshots() LayoutSnapshotRepository
}
