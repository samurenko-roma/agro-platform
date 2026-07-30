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

func (obj *ProductionUnit) AddDimensions(dim *Dimensions) {
	obj.Properties.Dimensions = dim
	if dim.Width != nil && dim.Length != nil {
		w := *dim.Width
		l := *dim.Length
		obj.Area = (w * l) / 10000
	}
}

func (obj *ProductionUnit) Occupy() {
	obj.Status = Growing
	obj.UpdatedAt = time.Now()
	//obj.AddEvent(NewProductionUnitOccupied(obj.ID))
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
	obj.Properties.Metadata["schema"] = schema

}

func BuildCode(parentCode string, unitType ProductionUnitType, seq int) string {
	part := fmt.Sprintf("%s%02d", unitType, seq)

	if parentCode == "" {
		return part
	}

	return parentCode + "-" + part
}
