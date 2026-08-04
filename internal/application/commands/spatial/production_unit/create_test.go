package productionunit_test

import (
	"context"
	"testing"

	"github.com/samurenkoroma/agro-platform/internal/application/commands/response"
	productionunitCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/spatial/production_unit"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	spatialrepo "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	inmemspatial "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/inmemory/spatial"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func newHandler() (*productionunitCmd.Handler, spatialrepo.SpatialProvider) {
	p := inmemspatial.NewProvider()
	uow := &testutil.FakeUoW{Provider: p}
	return productionunitCmd.NewProductionUnitHandler(uow), p
}

func TestCreate_RootUnit_Succeeds(t *testing.T) {
	h, _ := newHandler()
	res, err := h.Create(testutil.OrgCtx(), &productionunitCmd.CreateCommand{
		Type: pu.Field,
		Name: "Северное поле",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if res == nil {
		t.Fatal("expected result")
	}
}

func TestCreate_FailsWithoutOrgID(t *testing.T) {
	h, _ := newHandler()
	_, err := h.Create(context.Background(), &productionunitCmd.CreateCommand{Type: pu.Field, Name: "X"})
	if err == nil {
		t.Fatal("expected error without organization_id")
	}
}

func TestCreate_ChildWithDisallowedType_ReturnsInvalidHierarchy(t *testing.T) {
	h, _ := newHandler()
	ctx := testutil.OrgCtx()

	fieldRes, err := h.Create(ctx, &productionunitCmd.CreateCommand{Type: pu.Field, Name: "Поле"})
	if err != nil {
		t.Fatalf("create field: %v", err)
	}
	fieldID := fieldRes.(response.IdResponse).ID

	_, err = h.Create(ctx, &productionunitCmd.CreateCommand{
		Type:     pu.Slot, // FIELD -> SLOT недопустимо
		Name:     "Слот",
		ParentID: &fieldID,
	})
	if err != pu.ErrInvalidHierarchy {
		t.Fatalf("expected ErrInvalidHierarchy, got %v", err)
	}
}

func TestCreate_ChildWithAllowedType_Succeeds(t *testing.T) {
	h, _ := newHandler()
	ctx := testutil.OrgCtx()

	fieldRes, _ := h.Create(ctx, &productionunitCmd.CreateCommand{Type: pu.Field, Name: "Поле"})
	fieldID := fieldRes.(response.IdResponse).ID

	_, err := h.Create(ctx, &productionunitCmd.CreateCommand{
		Type:     pu.Block,
		Name:     "Блок 1",
		ParentID: &fieldID,
	})
	if err != nil {
		t.Fatalf("create block: %v", err)
	}
}

func TestCreate_ParentFromDifferentOrg_ReturnsNotFound(t *testing.T) {
	h, _ := newHandler()

	fieldRes, _ := h.Create(testutil.OrgCtx(), &productionunitCmd.CreateCommand{Type: pu.Field, Name: "Чужое поле"})
	fieldID := fieldRes.(response.IdResponse).ID

	// testutil.OrgCtx() каждый раз генерирует НОВУЮ случайную организацию —
	// это уже другая организация, чем та, что создала fieldID.
	_, err := h.Create(testutil.OrgCtx(), &productionunitCmd.CreateCommand{
		Type:     pu.Block,
		Name:     "Блок",
		ParentID: &fieldID,
	})
	if err != pu.ErrParentNotFound {
		t.Fatalf("expected ErrParentNotFound (межтенантная изоляция), got %v", err)
	}
}

func TestCreate_ManualAreaOverridesDimensions(t *testing.T) {
	h, p := newHandler()
	ctx := testutil.OrgCtx()

	width, length := 3.0, 0.8
	res, err := h.Create(ctx, &productionunitCmd.CreateCommand{
		Type:       pu.Field,
		Name:       "Поле",
		Dimensions: &pu.Dimensions{Width: &width, Length: &length},
		Area:       &pu.AreaInput{Value: 50, Unit: pu.Hectares},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	id := res.(response.IdResponse).ID

	orgID, _ := ctx.Value("organization_id").(string)
	unit, _ := p.Units().GetByID(ctx, id, vo.ID(orgID))
	if unit.Area != 500000 {
		t.Errorf("area: got %v, want 500000 (ручной ввод должен побеждать над dimensions)", unit.Area)
	}
}
