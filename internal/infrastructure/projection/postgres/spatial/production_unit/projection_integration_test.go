//go:build integration

package productionunit_test

import (
	"context"
	"testing"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	projection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/spatial/production_unit"
	pgspatial "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/spatial/production_unit"
	dbtest "github.com/samurenkoroma/agro-platform/internal/testutil/postgres"
)

var ctx = context.Background()

// Регрессия: раньше Tree() строил список "корней" поиском ParentID==nil
// в подвыборке CTE — для запрошенного НЕ-корневого узла (у которого в
// выборке есть свой родитель) список оставался пустым, и roots[0]
// паниковал (index out of range).
func TestGet_NonRootUnit_DoesNotPanic(t *testing.T) {
	db := dbtest.NewTestDB(t, "spatial")
	repo := pgspatial.New(db.Pool)
	proj := projection.New(db.Pool)

	ownerID := vo.NewID()
	fieldName := "Поле"
	field := pu.New(ownerID, nil, pu.Field, "FIELD01", &fieldName, 1)
	if err := repo.Save(ctx, field); err != nil {
		t.Fatalf("save field: %v", err)
	}

	blockName := "Блок"
	block := pu.New(ownerID, &field.ID, pu.Block, "FIELD01-BLOCK01", &blockName, 1)
	if err := repo.Save(ctx, block); err != nil {
		t.Fatalf("save block: %v", err)
	}

	dto, err := proj.Get(ctx, block.ID)
	if err != nil {
		t.Fatalf("get non-root unit (должно работать, не паниковать): %v", err)
	}
	if dto.ID != block.ID {
		t.Errorf("expected returned tree root to be the requested block itself, got %s", dto.ID)
	}
}

func TestGet_NonExistentUnit_ReturnsErrorNotPanic(t *testing.T) {
	db := dbtest.NewTestDB(t, "spatial")
	proj := projection.New(db.Pool)

	_, err := proj.Get(ctx, vo.NewID())
	if err != projection.ErrProductionUnitNotFound {
		t.Fatalf("expected ErrProductionUnitNotFound, got %v", err)
	}
}

func TestListRoots_UnitWithoutChildren_ChildrenIsEmptySliceNotNil(t *testing.T) {
	db := dbtest.NewTestDB(t, "spatial")
	repo := pgspatial.New(db.Pool)
	proj := projection.New(db.Pool)

	ownerID := vo.NewID()
	name := "Теплица"
	gh := pu.New(ownerID, nil, pu.Greenhouse, "GH01", &name, 1)
	if err := repo.Save(ctx, gh); err != nil {
		t.Fatalf("save: %v", err)
	}

	roots, err := proj.ListRoots(ctx, ownerID)
	if err != nil {
		t.Fatalf("list roots: %v", err)
	}
	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}
	if roots[0].Children == nil {
		t.Error("Children should be an empty slice, not nil")
	}
}
