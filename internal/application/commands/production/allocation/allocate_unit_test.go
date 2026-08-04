package allocation_test

import (
	"context"
	"testing"
	"time"

	allocationCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/allocation"
	growingcycleCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/growing_cycle"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/response"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	production "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	inmemproduction "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/inmemory/production"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func newHandlers() (*allocationCmd.Handler, *growingcycleCmd.Handler, production.ProductionProvider) {
	p := inmemproduction.NewProvider()
	uow := &testutil.FakeUoW{Provider: p}
	return allocationCmd.NewAllocationHandler(uow), growingcycleCmd.NewGrowingCycleHandler(uow), p
}

func timePtr(t time.Time) *time.Time { return &t }

func TestAllocateProductionUnit_FailsWithoutOrgID(t *testing.T) {
	h, _, _ := newHandlers()
	_, err := h.AllocateProductionUnit(context.Background(), &allocationCmd.AllocateProductionUnitCommand{
		CycleID: vo.NewID(), ProductionUnitID: vo.NewID(), Area: 10,
	})
	if err == nil {
		t.Fatal("expected error without organization_id")
	}
}

func TestAllocateProductionUnit_CycleFromDifferentOrg_Fails(t *testing.T) {
	h, cycleH, _ := newHandlers()

	cycleRes, _ := cycleH.Create(testutil.OrgCtx(), &growingcycleCmd.CreateCommand{
		CropID: vo.NewID(), Name: "Томат", Code: "TOM-1", Method: growingcycle.ProductionMethodSeedling,
	})
	cycleID := cycleRes.(response.IdResponse).ID

	// Другая (новая случайная) организация пытается аллоцировать чужой цикл.
	_, err := h.AllocateProductionUnit(testutil.OrgCtx(), &allocationCmd.AllocateProductionUnitCommand{
		CycleID: cycleID, ProductionUnitID: vo.NewID(), Area: 10,
	})
	if err != allocationCmd.ErrGrowingCycleNotFound {
		t.Fatalf("expected ErrGrowingCycleNotFound (межтенантная изоляция), got %v", err)
	}
}

func TestAllocateProductionUnit_Succeeds(t *testing.T) {
	h, cycleH, p := newHandlers()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	cycleRes, _ := cycleH.Create(ctx, &growingcycleCmd.CreateCommand{
		CropID: vo.NewID(), Name: "Томат", Code: "TOM-1", Method: growingcycle.ProductionMethodSeedling,
	})
	cycleID := cycleRes.(response.IdResponse).ID
	unitID := vo.NewID()

	res, err := h.AllocateProductionUnit(ctx, &allocationCmd.AllocateProductionUnitCommand{
		CycleID: cycleID, ProductionUnitID: unitID, Area: 10, StartedAt: timePtr(time.Now()),
	})
	if err != nil {
		t.Fatalf("allocate: %v", err)
	}
	allocID := res.(response.IdResponse).ID

	alloc, err := p.Allocation().GetByID(ctx, allocID, farmID)
	if err != nil || alloc == nil {
		t.Fatalf("allocation should be saved: %v", err)
	}
	if alloc.FarmID != farmID {
		t.Errorf("allocation FarmID: got %s, want %s", alloc.FarmID, farmID)
	}
}
