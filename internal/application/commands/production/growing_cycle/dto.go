package growingcycle

import (
	"time"

	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type CreateCommand struct {
	CropID    vo.ID  `json:"cropID" validate:"required"`
	VarietyID *vo.ID `json:"varietyID,omitempty"`

	Name       string                        `json:"name" validate:"required"`
	Code       string                        `json:"code" validate:"required"`
	Method     growingcycle.ProductionMethod `json:"method" validate:"required"`
	ProtocolID *vo.ID                        `json:"protocolID,omitempty"`
}

type StartGrowingCycleCMD struct {
	Name   string `json:"name" validate:"required"`
	Code   string `json:"code" validate:"required"`
	CropID vo.ID  `json:"cropID" validate:"required"`

	VarietyID  *vo.ID                        `json:"varietyID"`
	ProtocolID *vo.ID                        `json:"protocolID"`
	Stage      growingcycle.CycleStage       `json:"stage" validate:"omitempty,oneof=planning germination seedling vegetative flowering fruiting harvesting completed"`
	Method     growingcycle.ProductionMethod `json:"method" validate:"required"`

	ExpectedHarvestAt *time.Time `json:"expectedHarvestAt"`

	// Обязательны: "запустить" без единого физического размещения — это
	// просто create_cycle. Если нужно зарегистрировать цикл без размещения,
	// используйте production.create_cycle.
	Allocations []AllocationDTO `json:"allocations" validate:"required,min=1,dive"`
	Plantings   []PlantingDTO   `json:"plantings,omitempty" validate:"omitempty,dive"`
}

type AllocationDTO struct {
	ProductionUnitID vo.ID     `json:"productionUnitID" validate:"required"`
	Area             float64   `json:"area" validate:"required,gt=0"`
	StartedAt        time.Time `json:"startedAt"`
}

type PlantingDTO struct {
	PlantedAt time.Time `json:"plantedAt"`
	Quantity  float64   `json:"quantity" validate:"gt=0"`
}
