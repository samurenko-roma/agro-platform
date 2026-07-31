package layoutsnapshot

import (
	"context"
	"time"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type UnitDTO struct {
	ID             vo.ID     `json:"id"`
	OriginalUnitID vo.ID     `json:"originalUnitId"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	ParentID       *vo.ID    `json:"parentId,omitempty"`
	Capabilities   []string  `json:"capabilities,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type SnapshotDTO struct {
	ID          vo.ID     `json:"id"`
	FarmID      vo.ID     `json:"farmId"`
	Version     int       `json:"version"`
	Description *string   `json:"description,omitempty"`
	CreatedBy   vo.ID     `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	Units       []UnitDTO `json:"units"`
}

type Projection interface {
	Get(ctx context.Context, id vo.ID) (*SnapshotDTO, error)
	GetLatest(ctx context.Context, farmID vo.ID) (*SnapshotDTO, error)
}
