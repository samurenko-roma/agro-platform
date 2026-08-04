package planting_test

import (
	"testing"
	"time"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

func TestNew_SetsAllFields(t *testing.T) {
	farmID := vo.NewID()
	cycleID := vo.NewID()
	plantedAt := time.Now()

	p := planting.New(farmID, cycleID, plantedAt, 15.5)

	if p.ID.IsZero() {
		t.Error("expected non-zero ID")
	}
	if p.FarmID != farmID {
		t.Errorf("farmID: got %s, want %s", p.FarmID, farmID)
	}
	if p.CycleID != cycleID {
		t.Errorf("cycleID: got %s, want %s", p.CycleID, cycleID)
	}
	if p.Quantity != 15.5 {
		t.Errorf("quantity: got %v, want 15.5", p.Quantity)
	}
}
