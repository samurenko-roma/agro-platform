package productionunit

import (
	"encoding/json"
	"fmt"
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type ProductionUnitStatus string

const (
	Growing     ProductionUnitStatus = "growing"
	Preparation ProductionUnitStatus = "preparation"
	Empty       ProductionUnitStatus = "empty"
)

type Element struct {
	Id     string             `json:"id"`
	Type   ProductionUnitType `json:"type"`
	X      float64            `json:"x"`
	Y      float64            `json:"y"`
	Width  float64            `json:"width"`
	Length float64            `json:"length"`
	Name   string             `json:"name"`
}

type LayoutSchema struct {
	Beds []Element `json:"beds"`
}

type ProductionUnit struct {
	ev.BaseAggregate
	OwnerID    vo.ID
	ID         vo.ID
	ParentID   *vo.ID
	Type       ProductionUnitType
	Code       string
	Area       float64
	Geometry   vo.Geometry
	Properties *Properties
	Status     ProductionUnitStatus
	Sequence   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
}

func New(
	ownerID vo.ID,
	ParentId *vo.ID,
	unitType ProductionUnitType,
	code string,
	name *string,
	sequence int,
) *ProductionUnit {
	now := time.Now()

	displayName := ""
	if name != nil {
		displayName = *name
	}

	root := &ProductionUnit{
		ID:         vo.NewID(),
		ParentID:   ParentId,
		OwnerID:    ownerID,
		Code:       code,
		Sequence:   sequence,
		Type:       unitType,
		Status:     Empty,
		Properties: NewProps(displayName, ""),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	root.AddEvent(NewProductionUnitCreated(root.ID))

	return root
}

// AddDimensions задаёт размеры узла. Width/Length — в метрах.
// Площадь = width * length (м²).
func (obj *ProductionUnit) AddDimensions(dim *Dimensions) {
	obj.Properties.Dimensions = dim
	if dim.Width != nil && dim.Length != nil {
		obj.Area = *dim.Width * *dim.Length
	}
	obj.UpdatedAt = time.Now()
}

// SetGeometry задаёт геопривязку контура узла (актуально для полевых типов:
// FIELD, PLOT, BLOCK, BED — контур на карте фермы).
func (obj *ProductionUnit) SetGeometry(geometry vo.Geometry) {
	obj.Geometry = geometry
	obj.UpdatedAt = time.Now()
}

func (obj *ProductionUnit) Occupy() {
	obj.Status = Growing
	obj.UpdatedAt = time.Now()
	obj.AddEvent(NewProductionUnitOccupied(obj.ID))
}

func (obj *ProductionUnit) Release() {
	obj.Status = Empty
	obj.UpdatedAt = time.Now()
	obj.AddEvent(NewProductionUnitReleased(obj.ID))
}

func (obj *ProductionUnit) SetPreparation() {
	obj.Status = Preparation
	obj.UpdatedAt = time.Now()
	obj.AddEvent(NewProductionUnitInPreparation(obj.ID))
}

func (obj *ProductionUnit) UpdateSchema(schema json.RawMessage) {
	if obj.Properties == nil {
		obj.Properties = NewProps("", "")
	}
	obj.Properties.Metadata["schema"] = schema
	obj.UpdatedAt = time.Now()
}

// Archive помечает узел архивным. Архивные узлы:
//   - не отдаются в ListRoots/Tree (см. проекцию),
//   - не могут стать родителем нового узла (см. Create),
//   - не могут быть архивированы повторно.
func (obj *ProductionUnit) Archive() error {
	if obj.ArchivedAt != nil {
		return ErrAlreadyArchived
	}
	now := time.Now()
	obj.ArchivedAt = &now
	obj.UpdatedAt = now
	obj.AddEvent(NewProductionUnitArchived(obj.ID))
	return nil
}

func BuildCode(parentCode string, unitType ProductionUnitType, seq int) string {
	part := fmt.Sprintf("%s%02d", unitType, seq)

	if parentCode == "" {
		return part
	}

	return parentCode + "-" + part
}

// SetArea задаёт площадь явно (ручной ввод пользователя). Перекрывает то,
// что было бы вычислено из Dimensions — это осознанный override.
func (obj *ProductionUnit) SetArea(squareMeters float64) {
	obj.Area = squareMeters
	obj.UpdatedAt = time.Now()
}
