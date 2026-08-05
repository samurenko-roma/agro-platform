package repository

import "github.com/samurenkoroma/agro-platform/internal/shared/repository"

type ProductionProvider interface {
	repository.RepositoryProvider
	GrowingCycles() GrowingCycleRepository
	Harvests() HarvestBatchRepository
	Planting() PlantingRepository
	Allocation() AllocationRepository
}
