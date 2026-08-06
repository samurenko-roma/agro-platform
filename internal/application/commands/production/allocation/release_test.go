package allocation_test

import (
	"testing"
	"time"

	allocationCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/allocation"
	growingcycleCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/growing_cycle"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func createAllocation(t *testing.T) (*allocationCmd.Handler, vo.ID, vo.ID) {
	t.Helper()
	h, cycleH, _ := newHandlers()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	cycleRes, err := cycleH.Create(ctx, &growingcycleCmd.CreateCommand{
		CropID: vo.NewID(), Name: "Огурец", Code: "CUC-1", Method: growingcycle.ProductionMethodDirectSow,
	})
	if err != nil {
		t.Fatalf("create cycle: %v", err)
	}
	cycleID := cycleRes.(response.IdResponse).ID

	allocRes, err := h.AllocateProductionUnit(ctx, &allocationCmd.AllocateProductionUnitCommand{
		CycleID: cycleID, ProductionUnitID: vo.NewID(), Area: 5, StartedAt: timePtr(time.Now()),
	})
	if err != nil {
		t.Fatalf("allocate: %v", err)
	}
	return h, allocRes.(response.IdResponse).ID, farmID
}

func TestRelease_Succeeds(t *testing.T) {
	h, allocID, farmID := createAllocation(t)
	ctx := testutil.OrgCtxWithID(farmID)

	_, err := h.Release(ctx, &allocationCmd.ReleaseAllocationCommand{ID: allocID})
	if err != nil {
		t.Fatalf("release: %v", err)
	}
}

// Регрессия: раньше Release мутировал поле напрямую и не эмитил событие
// AllocationReleased — теперь через домен, событие есть (см. domain-тест
// TestRelease_SetsEndedAtAndEmitsEvent).
func TestRelease_Twice_Fails(t *testing.T) {
	h, allocID, farmID := createAllocation(t)
	ctx := testutil.OrgCtxWithID(farmID)

	h.Release(ctx, &allocationCmd.ReleaseAllocationCommand{ID: allocID})
	_, err := h.Release(ctx, &allocationCmd.ReleaseAllocationCommand{ID: allocID})
	if err == nil {
		t.Fatal("expected error releasing already-released allocation")
	}
}

func TestRelease_FromDifferentOrg_Fails(t *testing.T) {
	h, allocID, _ := createAllocation(t)

	_, err := h.Release(testutil.OrgCtx(), &allocationCmd.ReleaseAllocationCommand{ID: allocID})
	if err != allocationCmd.ErrAllocationNotFound {
		t.Fatalf("expected ErrAllocationNotFound (межтенантная изоляция), got %v", err)
	}
}
