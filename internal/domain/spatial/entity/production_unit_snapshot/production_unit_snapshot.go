package productionunitsnapshot

import (
	"time"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
)

type ProductionUnitSnapshot struct {
	ID vo.ID `json:"id"`

	SnapshotID vo.ID `json:"snapshotId"`

	OriginalUnitID vo.ID `json:"originalUnitId"`

	Type pu.ProductionUnitType `json:"type"`

	Name string `json:"name"`

	ParentID *vo.ID `json:"parentId,omitempty"`

	Capabilities []pu.Capability `json:"capabilities,omitempty"`

	Metadata vo.Metadata `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}
