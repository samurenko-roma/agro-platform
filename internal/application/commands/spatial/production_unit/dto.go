package productionunit

import (
	"encoding/json"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

type CreateCommand struct {
	Type         pu.ProductionUnitType `json:"type" validate:"required"`
	ParentID     *vo.ID                `json:"parentId,omitempty"`
	Capabilities []string              `json:"capabilities,omitempty" validate:"omitempty,oneof=soil irrigation fertigation drainage hydroponic aeroponic nutrientControl lighting climateControl sensorSupport automation slotBased mobile"`
	Name         string                `json:"name" validate:"required"`
	Dimensions   *pu.Dimensions        `json:"dimensions,omitempty"`
	Geometry     *vo.Geometry          `json:"geometry,omitempty"`
	Area         *pu.AreaInput         `json:"area,omitempty"`
}

type UpdateCommand struct {
	Id     vo.ID           `json:"id"`
	Schema json.RawMessage `json:"schema,omitempty"`
}

type ConfigureCommand struct {
	Id     vo.ID           `json:"id"`
	Schema pu.LayoutSchema `json:"schema"`
}

type ArchiveCommand struct {
	Id vo.ID `json:"id" validate:"required"`
}
