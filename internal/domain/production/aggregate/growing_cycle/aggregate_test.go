package growingcycle_test

import (
	"testing"

	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

var (
	farmID = vo.NewID()
	cropID = vo.NewID()
)

func newCycle(t *testing.T) *growingcycle.GrowingCycle {
	t.Helper()
	c := growingcycle.New(farmID, cropID, nil, nil, "Томат-1", "TOM-001", growingcycle.ProductionMethodSeedling)
	c.PullEvents()
	return c
}

func TestNew_DefaultsAreCorrect(t *testing.T) {
	c := growingcycle.New(farmID, cropID, nil, nil, "Томат-1", "TOM-001", growingcycle.ProductionMethodSeedling)

	if c.ID.IsZero() {
		t.Error("expected non-zero ID")
	}
	if c.FarmID != farmID {
		t.Errorf("farmID: got %s, want %s", c.FarmID, farmID)
	}
	if c.Status != growingcycle.StatusPlanned {
		t.Errorf("status: got %s, want StatusPlanned", c.Status)
	}
	if c.Stage != growingcycle.StagePlanning {
		t.Errorf("stage: got %s, want StagePlanning", c.Stage)
	}
}

func TestNew_EmitsCycleCreatedEvent(t *testing.T) {
	c := growingcycle.New(farmID, cropID, nil, nil, "Томат-1", "TOM-001", growingcycle.ProductionMethodSeedling)
	events := c.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType() != growingcycle.EventCreated {
		t.Errorf("event type: got %s, want %s", events[0].EventType(), growingcycle.EventCreated)
	}
}

func TestStart_SetsStatusActiveAndEmitsEvent(t *testing.T) {
	c := newCycle(t)
	c.Start()

	if c.Status != growingcycle.StatusActive {
		t.Errorf("status: got %s, want StatusActive", c.Status)
	}
	events := c.PullEvents()
	if len(events) != 1 || events[0].EventType() != growingcycle.EventStarted {
		t.Fatalf("expected 1 EventStarted, got %v", events)
	}
}

func TestChangeState_EmptyResetsToPlanning(t *testing.T) {
	c := newCycle(t)
	c.ChangeState(growingcycle.StageFlowering)
	c.ChangeState("")

	if c.Stage != growingcycle.StagePlanning {
		t.Errorf("stage: got %s, want StagePlanning", c.Stage)
	}
}

func TestChangeState_SetsGivenStage(t *testing.T) {
	c := newCycle(t)
	c.ChangeState(growingcycle.StageSeedling)

	if c.Stage != growingcycle.StageSeedling {
		t.Errorf("stage: got %s, want StageSeedling", c.Stage)
	}
}
