//go:build integration

package production_test

import (
	"context"
	"testing"
	"time"

	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	harvestbatch "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pgproduction "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/production"
	dbtest "github.com/samurenkoroma/agro-platform/internal/testutil/postgres"
)

var ctx = context.Background()

func mustSaveCycle(t *testing.T, db *dbtest.TestDB, farmID vo.ID) vo.ID {
	t.Helper()
	repo := pgproduction.NewGrowingCycleRepository(db.Pool)
	c := growingcycle.New(farmID, vo.NewID(), nil, nil, "Тест", "TEST-"+vo.NewID().String()[:8], growingcycle.ProductionMethodSeedling)
	if err := repo.Save(ctx, c); err != nil {
		t.Fatalf("save cycle: %v", err)
	}
	return c.ID
}

// Регрессия: раньше Save/ListByCycleID/Delete харвеста писали и читали
// production_plantings вместо production_harvest_batch — GetByID
// смотрел в правильную таблицу, которая из-за этого была вечно пуста.
func TestHarvestRepo_Postgres_SaveAndGetByID_RoundTrip(t *testing.T) {
	db := dbtest.NewTestDB(t, "production")
	repo := pgproduction.NewHarvestRepository(db.Pool)

	farmID := vo.NewID()
	cycleID := mustSaveCycle(t, db, farmID)

	batch := harvestbatch.New(farmID, cycleID, time.Now(), 320.5)
	if err := repo.Save(ctx, batch); err != nil {
		t.Fatalf("save harvest: %v", err)
	}

	got, err := repo.GetByID(ctx, batch.ID, farmID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected harvest batch to be found after save (регрессия: Save/GetByID раньше смотрели в разные таблицы)")
	}
	if got.Quantity != 320.5 {
		t.Errorf("quantity: got %v, want 320.5", got.Quantity)
	}
}

func TestHarvestRepo_Postgres_DoesNotPolluteProductionPlantings(t *testing.T) {
	db := dbtest.NewTestDB(t, "production")
	harvestRepo := pgproduction.NewHarvestRepository(db.Pool)
	plantingRepo := pgproduction.NewPlantingRepository(db.Pool)

	farmID := vo.NewID()
	cycleID := mustSaveCycle(t, db, farmID)

	batch := harvestbatch.New(farmID, cycleID, time.Now(), 100)
	if err := harvestRepo.Save(ctx, batch); err != nil {
		t.Fatalf("save harvest: %v", err)
	}

	plantings, err := plantingRepo.ListByCycleID(ctx, cycleID)
	if err != nil {
		t.Fatalf("list plantings: %v", err)
	}
	if len(plantings) != 0 {
		t.Errorf("expected 0 plantings (харвест не должен утекать в таблицу посадок), got %d", len(plantings))
	}
}
