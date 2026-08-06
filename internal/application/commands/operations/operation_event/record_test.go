package operationevent_test

import (
	"context"
	"testing"

	operationevent "github.com/samurenkoroma/agro-platform/internal/application/commands/operations/operation_event"
	opsrepo "github.com/samurenkoroma/agro-platform/internal/domain/operations/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	inmemops "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/inmemory/operations"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func newHandler() (*operationevent.Handler, opsrepo.OperationsProvider) {
	p := inmemops.NewProvider()
	uow := &testutil.FakeUoW{Provider: p}
	return operationevent.NewOperationHandler(uow), p
}

// --- Record ---

func TestRecord_ReturnsID(t *testing.T) {
	h, _ := newHandler()
	result, err := h.Record(testutil.UserCtx(), &operationevent.RecordOperationCommand{
		Type: "IRRIGATED",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestRecord_FailsWithoutOrgID(t *testing.T) {
	h, _ := newHandler()
	_, err := h.Record(context.Background(), &operationevent.RecordOperationCommand{Type: "IRRIGATED"})
	if err == nil {
		t.Fatal("expected error when organization_id is missing")
	}
}

func TestRecord_FailsWithWrongPayloadType(t *testing.T) {
	h, _ := newHandler()
	_, err := h.Record(testutil.UserCtx(), "bad payload")
	if err == nil {
		t.Fatal("expected error for invalid payload type")
	}
}

func TestRecord_PersistsOperationEvent(t *testing.T) {
	h, p := newHandler()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	h.Record(ctx, &operationevent.RecordOperationCommand{
		Type:    "FERTILIZED",
		Payload: map[string]any{"amount": 500},
	})

	ops, err := p.Operations().List(ctx, opsrepo.OperationFilter{FarmID: farmID})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(ops))
	}
	if ops[0].Type != "FERTILIZED" {
		t.Errorf("type: got %s, want FERTILIZED", ops[0].Type)
	}
}

func TestRecord_PayloadIsStored(t *testing.T) {
	h, p := newHandler()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	h.Record(ctx, &operationevent.RecordOperationCommand{
		Type:    "IRRIGATED",
		Payload: map[string]any{"volumeLiters": 100.0, "method": "drip"},
	})

	ops, _ := p.Operations().List(ctx, opsrepo.OperationFilter{FarmID: farmID})
	if ops[0].Payload["volumeLiters"] != 100.0 {
		t.Errorf("payload volumeLiters: got %v, want 100", ops[0].Payload["volumeLiters"])
	}
	if ops[0].Payload["method"] != "drip" {
		t.Errorf("payload method: got %v, want drip", ops[0].Payload["method"])
	}
}

func TestRecord_AutoCreatesTimeline(t *testing.T) {
	h, p := newHandler()
	farmID := vo.NewID()
	cycleID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)
	cycleStr := cycleID.String()

	h.Record(ctx, &operationevent.RecordOperationCommand{
		Type:           "IRRIGATED",
		GrowingCycleID: &cycleStr,
	})

	tl, err := p.Timelines().GetByOwner(ctx, farmID, &cycleID)
	if err != nil {
		t.Fatalf("get timeline: %v", err)
	}
	if tl == nil {
		t.Fatal("expected timeline to be created automatically")
	}
	if len(tl.Items) != 1 {
		t.Fatalf("expected 1 timeline item, got %d", len(tl.Items))
	}
}

func TestRecord_AppendsToExistingTimeline(t *testing.T) {
	h, p := newHandler()
	farmID := vo.NewID()
	cycleID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)
	cycleStr := cycleID.String()

	h.Record(ctx, &operationevent.RecordOperationCommand{Type: "IRRIGATED", GrowingCycleID: &cycleStr})
	h.Record(ctx, &operationevent.RecordOperationCommand{Type: "FERTILIZED", GrowingCycleID: &cycleStr})
	h.Record(ctx, &operationevent.RecordOperationCommand{Type: "PH_ADJUSTED", GrowingCycleID: &cycleStr})

	tl, _ := p.Timelines().GetByOwner(ctx, farmID, &cycleID)
	if len(tl.Items) != 3 {
		t.Fatalf("expected 3 timeline items, got %d", len(tl.Items))
	}
}

func TestRecord_TimelineItemHasCorrectTitle(t *testing.T) {
	h, p := newHandler()
	farmID := vo.NewID()
	cycleID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)
	cycleStr := cycleID.String()

	h.Record(ctx, &operationevent.RecordOperationCommand{Type: "FERTILIZED", GrowingCycleID: &cycleStr})

	tl, _ := p.Timelines().GetByOwner(ctx, farmID, &cycleID)
	if tl.Items[0].Title != "FERTILIZED" {
		t.Errorf("item title: got %s, want FERTILIZED", tl.Items[0].Title)
	}
}

func TestRecord_WithPerformedBy(t *testing.T) {
	h, p := newHandler()
	farmID := vo.NewID()
	userID := vo.NewID()
	ctx := context.WithValue(testutil.OrgCtxWithID(farmID), "user_id", userID.String())

	h.Record(ctx, &operationevent.RecordOperationCommand{Type: "IRRIGATED"})

	ops, _ := p.Operations().List(ctx, opsrepo.OperationFilter{FarmID: farmID})
	if ops[0].PerformedBy == nil {
		t.Fatal("PerformedBy should be set from context user_id")
	}
	if *ops[0].PerformedBy != userID {
		t.Errorf("PerformedBy: got %s, want %s", *ops[0].PerformedBy, userID)
	}
}

func TestRecord_MultipleOperationsInDifferentOrgs(t *testing.T) {
	h, p := newHandler()
	farm1 := vo.NewID()
	farm2 := vo.NewID()

	h.Record(testutil.OrgCtxWithID(farm1), &operationevent.RecordOperationCommand{Type: "IRRIGATED"})
	h.Record(testutil.OrgCtxWithID(farm2), &operationevent.RecordOperationCommand{Type: "FERTILIZED"})

	ops1, _ := p.Operations().List(context.Background(), opsrepo.OperationFilter{FarmID: farm1})
	ops2, _ := p.Operations().List(context.Background(), opsrepo.OperationFilter{FarmID: farm2})

	if len(ops1) != 1 {
		t.Errorf("farm1: expected 1 operation, got %d", len(ops1))
	}
	if len(ops2) != 1 {
		t.Errorf("farm2: expected 1 operation, got %d", len(ops2))
	}
}

// подавляем warning о неиспользуемом импорте
var _ = response.IdResponse{}
