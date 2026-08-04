package allocation_test

import (
	"testing"
	"time"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

var (
	farmID    = vo.NewID()
	cycleID   = vo.NewID()
	unitID    = vo.NewID()
	otherUnit = vo.NewID()
)

func newAllocation(t *testing.T) *allocation.Allocation {
	t.Helper()
	now := time.Now()
	a := allocation.New(farmID, cycleID, unitID, 12.5, &now)
	a.PullEvents()
	return a
}

func TestNew_EmitsAllocationAllocatedEvent(t *testing.T) {
	now := time.Now()
	a := allocation.New(farmID, cycleID, unitID, 12.5, &now)

	events := a.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(allocation.AllocationAllocated)
	if !ok {
		t.Fatalf("expected AllocationAllocated, got %T", events[0])
	}
	if evt.ProductionUnitID != unitID {
		t.Errorf("ProductionUnitID: got %s, want %s", evt.ProductionUnitID, unitID)
	}
}

func TestRelease_SetsEndedAtAndEmitsEvent(t *testing.T) {
	a := newAllocation(t)
	releasedAt := time.Now()

	if err := a.Release(releasedAt); err != nil {
		t.Fatalf("release: %v", err)
	}
	if a.EndedAt == nil || !a.EndedAt.Equal(releasedAt) {
		t.Errorf("EndedAt not set correctly")
	}
	events := a.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	evt, ok := events[0].(allocation.AllocationReleased)
	if !ok {
		t.Fatalf("expected AllocationReleased, got %T", events[0])
	}
	if evt.ProductionUnitID != unitID {
		t.Errorf("ProductionUnitID: got %s, want %s", evt.ProductionUnitID, unitID)
	}
}

func TestRelease_Twice_ReturnsError(t *testing.T) {
	a := newAllocation(t)
	a.Release(time.Now())
	if err := a.Release(time.Now()); err != allocation.ErrAlreadyReleased {
		t.Errorf("expected ErrAlreadyReleased, got %v", err)
	}
}

// --- Reallocate ---
// Регрессионные тесты на баг: раньше change_allocation мутировал поле
// напрямую и не эмитил ничего — Spatial никогда не узнавал о переносе.

func TestReallocate_SameUnit_NoEventsEmitted(t *testing.T) {
	a := newAllocation(t)
	now := time.Now()
	a.Reallocate(unitID, 20, &now, nil)

	if a.Area != 20 {
		t.Errorf("area: got %v, want 20", a.Area)
	}
	events := a.PullEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 events when unit unchanged, got %d", len(events))
	}
}

func TestReallocate_DifferentUnit_EmitsReleaseThenAllocate(t *testing.T) {
	a := newAllocation(t)
	now := time.Now()
	a.Reallocate(otherUnit, 20, &now, nil)

	if a.ProductionUnitID != otherUnit {
		t.Errorf("ProductionUnitID: got %s, want %s", a.ProductionUnitID, otherUnit)
	}

	events := a.PullEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	released, ok := events[0].(allocation.AllocationReleased)
	if !ok || released.ProductionUnitID != unitID {
		t.Errorf("expected first event to release OLD unit %s, got %+v", unitID, events[0])
	}
	allocated, ok := events[1].(allocation.AllocationAllocated)
	if !ok || allocated.ProductionUnitID != otherUnit {
		t.Errorf("expected second event to allocate NEW unit %s, got %+v", otherUnit, events[1])
	}
}
