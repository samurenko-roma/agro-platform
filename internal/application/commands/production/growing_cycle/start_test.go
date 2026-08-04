package growingcycle_test

import (
	"testing"
	"time"

	growingcycleCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/production/growing_cycle"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/response"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	production "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	inmemproduction "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/inmemory/production"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func newStartHandler() (*growingcycleCmd.Handler, production.ProductionProvider) {
	p := inmemproduction.NewProvider()
	uow := &testutil.FakeUoW{Provider: p}
	return growingcycleCmd.NewGrowingCycleHandler(uow), p
}

func TestStart_ForcesStatusActive(t *testing.T) {
	h, p := newStartHandler()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	res, err := h.Start(ctx, &growingcycleCmd.StartGrowingCycleCMD{
		CropID: vo.NewID(), Name: "Томат", Code: "TOM-1", Method: growingcycle.ProductionMethodSeedling,
		Allocations: []growingcycleCmd.AllocationDTO{
			{ProductionUnitID: vo.NewID(), Area: 5, StartedAt: time.Now()},
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	cycleID := res.(response.IdResponse).ID

	cycle, err := p.GrowingCycles().GetByID(ctx, cycleID, farmID)
	if err != nil || cycle == nil {
		t.Fatalf("cycle should exist: %v", err)
	}
	if cycle.Status != growingcycle.StatusActive {
		t.Errorf("status: got %s, want Active", cycle.Status)
	}
}

// Регрессия: раньше cmd.Plantings объявлялось в DTO и полностью
// игнорировалось в start.go — клиент думал, что регистрирует посевы,
// а они молча пропадали.
func TestStart_CreatesAllocationsAndPlantings(t *testing.T) {
	h, p := newStartHandler()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	res, err := h.Start(ctx, &growingcycleCmd.StartGrowingCycleCMD{
		CropID: vo.NewID(), Name: "Огурец", Code: "CUC-1", Method: growingcycle.ProductionMethodDirectSow,
		Allocations: []growingcycleCmd.AllocationDTO{
			{ProductionUnitID: vo.NewID(), Area: 5, StartedAt: time.Now()},
		},
		Plantings: []growingcycleCmd.PlantingDTO{
			{PlantedAt: time.Now(), Quantity: 250},
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	cycleID := res.(response.IdResponse).ID

	allocs, _ := p.Allocation().ListByCycleID(ctx, cycleID)
	if len(allocs) != 1 {
		t.Fatalf("expected 1 allocation, got %d", len(allocs))
	}

	plantings, _ := p.Planting().ListByCycleID(ctx, cycleID)
	if len(plantings) != 1 {
		t.Fatalf("expected 1 planting (Plantings должен использоваться, не игнорироваться), got %d", len(plantings))
	}
	if plantings[0].Quantity != 250 {
		t.Errorf("planting quantity: got %v, want 250", plantings[0].Quantity)
	}
}
