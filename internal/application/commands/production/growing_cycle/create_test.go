package growingcycle_test

import (
	"testing"

	growingcycleCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/growing_cycle"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func TestCreate_StaysInPlannedWithoutAllocations(t *testing.T) {
	h, p := newStartHandler()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	res, err := h.Create(ctx, &growingcycleCmd.CreateCommand{
		CropID: vo.NewID(), Name: "Перец", Code: "PEP-1", Method: growingcycle.ProductionMethodSeedling,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	cycleID := res.(response.IdResponse).ID

	cycle, _ := p.GrowingCycles().GetByID(ctx, cycleID, farmID)
	if cycle.Status != growingcycle.StatusPlanned {
		t.Errorf("status: got %s, want Planned (create_cycle не должен активировать)", cycle.Status)
	}

	allocs, _ := p.Allocation().ListByCycleID(ctx, cycleID)
	if len(allocs) != 0 {
		t.Errorf("expected 0 allocations from create_cycle, got %d", len(allocs))
	}
}
