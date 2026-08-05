//go:build integration

package allocation_test

import (
	"context"
	"testing"
	"time"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	projection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/production/allocation"
	pgproduction "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/production"
	pgspatial "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/spatial/production_unit"
	dbtest "github.com/samurenkoroma/agro-platform/internal/testutil/postgres"
)

var ctx = context.Background()

// Регрессия: раньше JOIN сравнивал a.production_unit_id = a.id (опечатка,
// обе стороны из одной таблицы) — unit.code всегда был NULL в Get()/List().
func TestAllocationProjection_Get_ReturnsProductionUnitCode(t *testing.T) {
	db := dbtest.NewTestDB(t, "production", "spatial")

	unitRepo := pgspatial.New(db.Pool)
	cycleRepo := pgproduction.NewGrowingCycleRepository(db.Pool)
	allocRepo := pgproduction.NewAllocationRepository(db.Pool)
	proj := projection.New(db.Pool)

	farmID := vo.NewID()
	name := "Грядка 1"
	unit := pu.New(farmID, nil, pu.Bed, "BED01", &name, 1)
	if err := unitRepo.Save(ctx, unit); err != nil {
		t.Fatalf("save unit: %v", err)
	}

	cycle := growingcycle.New(farmID, vo.NewID(), nil, nil, "Тест", "TEST-ALLOC-1", growingcycle.ProductionMethodSeedling)
	if err := cycleRepo.Save(ctx, cycle); err != nil {
		t.Fatalf("save cycle: %v", err)
	}

	now := time.Now()
	alloc := allocation.New(farmID, cycle.ID, unit.ID, 2.4, &now)
	if err := allocRepo.Save(ctx, alloc); err != nil {
		t.Fatalf("save allocation: %v", err)
	}

	dto, err := proj.Get(ctx, alloc.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if dto.ProductionUnitName != "BED01" {
		t.Errorf("ProductionUnitName: got %q, want %q (JOIN должен резолвить код юнита)", dto.ProductionUnitName, "BED01")
	}
}

// Регрессия: раньше List() фильтровал по unit.owner_id, который из-за
// сломанного JOIN всегда был NULL — список был вечно пуст независимо от
// реальных данных.
func TestAllocationProjection_List_ReturnsAllocationsForFarm(t *testing.T) {
	db := dbtest.NewTestDB(t, "production", "spatial")

	unitRepo := pgspatial.New(db.Pool)
	cycleRepo := pgproduction.NewGrowingCycleRepository(db.Pool)
	allocRepo := pgproduction.NewAllocationRepository(db.Pool)
	proj := projection.New(db.Pool)

	farmID := vo.NewID()
	name := "Грядка 1"
	unit := pu.New(farmID, nil, pu.Bed, "BED01", &name, 1)
	unitRepo.Save(ctx, unit)

	cycle := growingcycle.New(farmID, vo.NewID(), nil, nil, "Тест", "TEST-ALLOC-2", growingcycle.ProductionMethodSeedling)
	cycleRepo.Save(ctx, cycle)

	now := time.Now()
	alloc := allocation.New(farmID, cycle.ID, unit.ID, 2.4, &now)
	allocRepo.Save(ctx, alloc)

	list, err := proj.List(ctx, farmID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 allocation, got %d (раньше List() всегда возвращал пустой список)", len(list))
	}
}
