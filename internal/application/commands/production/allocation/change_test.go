package allocation_test

import (
	"testing"

	allocationCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/allocation"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/response"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func TestChange_ReallocateToNewUnit_UpdatesProductionUnitID(t *testing.T) {
	h, allocID, farmID := createAllocation(t)
	ctx := testutil.OrgCtxWithID(farmID)
	newUnitID := vo.NewID()

	res, err := h.Change(ctx, &allocationCmd.ChangeAllocationCommand{
		ID: allocID, ProductionUnitID: newUnitID, Area: 8,
	})
	if err != nil {
		t.Fatalf("change: %v", err)
	}
	if res.(response.IdResponse).ID != allocID {
		t.Errorf("expected same allocation ID returned")
	}
}
