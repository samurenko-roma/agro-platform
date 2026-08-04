package productionunit_test

import (
	"encoding/json"
	"testing"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

var ownerID = vo.NewID()

func newUnit(t *testing.T) *pu.ProductionUnit {
	t.Helper()
	name := "Грядка 1"
	unit := pu.New(ownerID, nil, pu.Bed, "BED01", &name, 1)
	unit.PullEvents()
	return unit
}

// --- New ---

func TestNew_DefaultsAreCorrect(t *testing.T) {
	name := "Поле"
	unit := pu.New(ownerID, nil, pu.Field, "FIELD01", &name, 1)

	if unit.ID.IsZero() {
		t.Error("expected non-zero ID")
	}
	if unit.OwnerID != ownerID {
		t.Errorf("ownerID: got %s, want %s", unit.OwnerID, ownerID)
	}
	if unit.Status != pu.Empty {
		t.Errorf("status: got %s, want Empty", unit.Status)
	}
	if unit.Properties == nil {
		t.Fatal("Properties should be initialized")
	}
	if unit.Properties.Metadata["name"] != "Поле" {
		t.Errorf("name in metadata: got %v", unit.Properties.Metadata["name"])
	}
}

func TestNew_NilNameDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("New panicked with nil name: %v", r)
		}
	}()
	unit := pu.New(ownerID, nil, pu.Field, "FIELD01", nil, 1)
	if unit.Properties.Metadata["name"] != "" {
		t.Errorf("expected empty name, got %v", unit.Properties.Metadata["name"])
	}
}

func TestNew_EmitsProductionUnitCreatedEvent(t *testing.T) {
	name := "Поле"
	unit := pu.New(ownerID, nil, pu.Field, "FIELD01", &name, 1)
	events := unit.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType() != pu.EventCreated {
		t.Errorf("event type: got %s, want %s", events[0].EventType(), pu.EventCreated)
	}
}

// --- AddDimensions ---

func TestAddDimensions_ComputesAreaInSquareMeters(t *testing.T) {
	unit := newUnit(t)
	width, length := 3.0, 0.8
	unit.AddDimensions(&pu.Dimensions{Width: &width, Length: &length})

	want := 2.4
	if unit.Area != want {
		t.Errorf("area: got %v, want %v", unit.Area, want)
	}
}

func TestAddDimensions_WithoutWidthOrLength_AreaUnchanged(t *testing.T) {
	unit := newUnit(t)
	length := 0.8
	unit.AddDimensions(&pu.Dimensions{Length: &length})
	if unit.Area != 0 {
		t.Errorf("expected area 0 without width, got %v", unit.Area)
	}
}

// --- SetArea ---

func TestSetArea_OverridesDimensionsArea(t *testing.T) {
	unit := newUnit(t)
	width, length := 3.0, 0.8
	unit.AddDimensions(&pu.Dimensions{Width: &width, Length: &length})
	unit.SetArea(500000) // 50 га в м²

	if unit.Area != 500000 {
		t.Errorf("area: got %v, want 500000 (ручной ввод должен побеждать)", unit.Area)
	}
}

// --- Occupy / Release / SetPreparation ---

func TestOccupy_ChangesStatusAndEmitsEvent(t *testing.T) {
	unit := newUnit(t)
	unit.Occupy()

	if unit.Status != pu.Growing {
		t.Errorf("status: got %s, want Growing", unit.Status)
	}
	events := unit.PullEvents()
	if len(events) != 1 || events[0].EventType() != pu.EventOccupied {
		t.Fatalf("expected 1 EventOccupied, got %v", events)
	}
}

func TestRelease_ChangesStatusAndEmitsEvent(t *testing.T) {
	unit := newUnit(t)
	unit.Occupy()
	unit.PullEvents()
	unit.Release()

	if unit.Status != pu.Empty {
		t.Errorf("status: got %s, want Empty", unit.Status)
	}
	events := unit.PullEvents()
	if len(events) != 1 || events[0].EventType() != pu.EventReleased {
		t.Fatalf("expected 1 EventReleased, got %v", events)
	}
}

func TestSetPreparation_ChangesStatusAndEmitsEvent(t *testing.T) {
	unit := newUnit(t)
	unit.SetPreparation()

	if unit.Status != pu.Preparation {
		t.Errorf("status: got %s, want Preparation", unit.Status)
	}
	events := unit.PullEvents()
	if len(events) != 1 || events[0].EventType() != pu.EventInPreparation {
		t.Fatalf("expected 1 EventInPreparation, got %v", events)
	}
}

// --- UpdateSchema ---

func TestUpdateSchema_StoresRawSchema(t *testing.T) {
	unit := newUnit(t)
	schema := json.RawMessage(`{"beds":[]}`)
	unit.UpdateSchema(schema)

	stored, ok := unit.Properties.Metadata["schema"].(json.RawMessage)
	if !ok || string(stored) != string(schema) {
		t.Errorf("schema not stored correctly: %v", unit.Properties.Metadata["schema"])
	}
}

func TestUpdateSchema_NilPropertiesDoesNotPanic(t *testing.T) {
	unit := &pu.ProductionUnit{} // Properties == nil, как если бы юнит был собран вручную/некорректно
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("UpdateSchema panicked on nil Properties: %v", r)
		}
	}()
	unit.UpdateSchema(json.RawMessage(`{}`))
}

// --- Archive ---

func TestArchive_SetsArchivedAtAndEmitsEvent(t *testing.T) {
	unit := newUnit(t)
	if err := unit.Archive(); err != nil {
		t.Fatalf("archive: %v", err)
	}
	if unit.ArchivedAt == nil {
		t.Fatal("ArchivedAt should be set")
	}
	events := unit.PullEvents()
	if len(events) != 1 || events[0].EventType() != pu.EventArchived {
		t.Fatalf("expected 1 EventArchived, got %v", events)
	}
}

func TestArchive_Twice_ReturnsError(t *testing.T) {
	unit := newUnit(t)
	unit.Archive()
	if err := unit.Archive(); err != pu.ErrAlreadyArchived {
		t.Errorf("expected ErrAlreadyArchived, got %v", err)
	}
}

// --- BuildCode ---

func TestBuildCode_WithParent(t *testing.T) {
	code := pu.BuildCode("FIELD01", pu.Block, 2)
	want := "FIELD01-BLOCK02"
	if code != want {
		t.Errorf("code: got %s, want %s", code, want)
	}
}

func TestBuildCode_WithoutParent(t *testing.T) {
	code := pu.BuildCode("", pu.Field, 1)
	want := "FIELD01"
	if code != want {
		t.Errorf("code: got %s, want %s", code, want)
	}
}
