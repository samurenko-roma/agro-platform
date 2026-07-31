package spatial

import (
	"context"
	"fmt"

	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	topology "github.com/samurenkoroma/agro-platform/internal/domain/spatial/service"
)

type LayoutGenerator interface {
	Generate(ctx context.Context, parent *pu.ProductionUnit, schema pu.LayoutSchema) error
}

type generator struct {
	repo     repository.ProductionUnitRepository
	exec     uow.Execution
	topology topology.TopologyRules
}

func (g generator) Generate(ctx context.Context, parent *pu.ProductionUnit, schema pu.LayoutSchema) error {
	nextSeq := make(map[pu.ProductionUnitType]int)

	for _, element := range schema.Beds {
		unitType := element.Type.ToUpper()

		if !g.topology.CanAttach(parent.Type, unitType) {
			return fmt.Errorf("layout: %s cannot contain %s", parent.Type, unitType)
		}

		seq, ok := nextSeq[unitType]
		if !ok {
			s, err := g.repo.GetNextSequence(ctx, parent.OwnerID, &parent.ID, unitType)
			if err != nil {
				return err
			}
			seq = s
		}

		code := pu.BuildCode(parent.Code, unitType, seq)

		name := element.Name
		unit := pu.New(parent.OwnerID, &parent.ID, unitType, code, &name, seq)
		unit.AddDimensions(&pu.Dimensions{Length: &element.Length, Width: &element.Width})
		unit.Properties.Position = &valueobject.Position{
			X: element.X,
			Y: element.Y,
		}

		if err := g.repo.Save(ctx, unit); err != nil {
			return err
		}
		g.exec.RegisterAggregate(unit)

		nextSeq[unitType] = seq + 1
	}

	return nil
}

func NewUnitLayoutGenerator(repo repository.ProductionUnitRepository, exec uow.Execution) LayoutGenerator {
	return generator{
		repo:     repo,
		exec:     exec,
		topology: topology.DefaultTopology{},
	}
}
