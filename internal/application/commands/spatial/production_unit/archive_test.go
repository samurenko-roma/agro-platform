package productionunit_test

import (
	"errors"
	"testing"

	productionunitCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/spatial/production_unit"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func TestArchive_UnitWithoutChildren_Succeeds(t *testing.T) {
	h, _ := newHandler()
	ctx := testutil.OrgCtx()

	res, _ := h.Create(ctx, &productionunitCmd.CreateCommand{Type: pu.Field, Name: "Поле"})
	id := res.(response.IdResponse).ID

	if _, err := h.Archive(ctx, &productionunitCmd.ArchiveCommand{Id: id}); err != nil {
		t.Fatalf("archive: %v", err)
	}
}

func TestArchive_UnitWithActiveChildren_Fails(t *testing.T) {
	h, _ := newHandler()
	ctx := testutil.OrgCtx()

	fieldRes, _ := h.Create(ctx, &productionunitCmd.CreateCommand{Type: pu.Field, Name: "Поле"})
	fieldID := fieldRes.(response.IdResponse).ID

	h.Create(ctx, &productionunitCmd.CreateCommand{Type: pu.Block, Name: "Блок", ParentID: &fieldID})

	_, err := h.Archive(ctx, &productionunitCmd.ArchiveCommand{Id: fieldID})
	if !errors.Is(err, productionunitCmd.ErrHasActiveChildren) {
		t.Fatalf("expected ErrHasActiveChildren, got %v", err)
	}
}

func TestArchive_AlreadyArchived_Fails(t *testing.T) {
	h, _ := newHandler()
	ctx := testutil.OrgCtx()

	res, _ := h.Create(ctx, &productionunitCmd.CreateCommand{Type: pu.Field, Name: "Поле"})
	id := res.(response.IdResponse).ID

	h.Archive(ctx, &productionunitCmd.ArchiveCommand{Id: id})
	_, err := h.Archive(ctx, &productionunitCmd.ArchiveCommand{Id: id})
	if !errors.Is(err, pu.ErrAlreadyArchived) {
		t.Fatalf("expected ErrAlreadyArchived, got %v", err)
	}
}
