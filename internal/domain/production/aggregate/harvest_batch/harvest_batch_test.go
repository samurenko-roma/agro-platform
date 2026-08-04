package harvestbatch_test

import (
	"testing"
	"time"

	harvestbatch "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

func TestNew_SetsAllFields(t *testing.T) {
	farmID := vo.NewID()
	cycleID := vo.NewID()
	harvestedAt := time.Now()

	h := harvestbatch.New(farmID, cycleID, harvestedAt, 320.0)

	if h.ID.IsZero() {
		t.Error("expected non-zero ID")
	}
	if h.FarmID != farmID {
		t.Errorf("farmID: got %s, want %s", h.FarmID, farmID)
	}
	if h.CycleID != cycleID {
		t.Errorf("cycleID: got %s, want %s", h.CycleID, cycleID)
	}
	if h.Quantity != 320.0 {
		t.Errorf("quantity: got %v, want 320", h.Quantity)
	}
}
