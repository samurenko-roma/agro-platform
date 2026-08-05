//go:build integration

package growingcycle_test

import (
	"context"
	"testing"

	query "github.com/samurenkoroma/agro-platform/internal/application/queries/production/growing_cycle"
	agronomycrop "github.com/samurenkoroma/agro-platform/internal/domain/agronomy/aggregate/crop"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	projection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/production/growing_cycle"
	pgagronomy "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/agronomy"
	pgproduction "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/production"
	dbtest "github.com/samurenkoroma/agro-platform/internal/testutil/postgres"
)

var ctx = context.Background()

func mustSaveCrop(t *testing.T, db *dbtest.TestDB) vo.ID {
	t.Helper()
	repo := pgagronomy.NewCropRepository(db.Pool)
	c := agronomycrop.New("Тестовая культура", agronomycrop.Vegetable, "Solanaceae", agronomycrop.AgronomyProfile{})
	if err := repo.Save(ctx, c); err != nil {
		t.Fatalf("save crop: %v", err)
	}
	return c.ID
}

// Регрессия: раньше *growingDays разыменовывался без nil-проверки — цикл
// без единой аллокации (самый обычный случай — только что создан) гарантированно
// ронял List() паникой (nil pointer dereference).
func TestGrowingCycleProjection_List_CycleWithoutAllocations_DoesNotPanic(t *testing.T) {
	db := dbtest.NewTestDB(t, "production", "agronomy")

	cycleRepo := pgproduction.NewGrowingCycleRepository(db.Pool)
	proj := projection.New(db.Pool)

	farmID := vo.NewID()
	cropID := mustSaveCrop(t, db)
	cycle := growingcycle.New(farmID, cropID, nil, nil, "Без размещения", "NOALLOC-1", growingcycle.ProductionMethodSeedling)

	if err := cycleRepo.Save(ctx, cycle); err != nil {
		t.Fatalf("save cycle: %v", err)
	}

	result, err := proj.List(ctx, query.FilterCycle{OwnerId: farmID})
	if err != nil {
		t.Fatalf("list (не должен паниковать/возвращать ошибку): %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 crop group, got %d", len(result))
	}
	if result[0].Progress != 0 {
		t.Errorf("progress без аллокаций должен быть 0, got %d", result[0].Progress)
	}
}
