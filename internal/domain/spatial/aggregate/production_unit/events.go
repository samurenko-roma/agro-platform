package productionunit

import (
	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/event"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

const (
	EventCreated = "production_unit.created"

	EventAttached = "production_unit.attached"

	EventOccupied      = "production_unit.occupied"
	EventReleased      = "production_unit.released"
	EventInPreparation = "production_unit.in_preparation"
	EventArchived      = "production_unit.archived"
)

type ProductionUnitOccupied struct{ ev.BaseEvent }

func NewProductionUnitOccupied(id vo.ID) ProductionUnitOccupied {
	return ProductionUnitOccupied{ev.NewBaseEvent(id, EventOccupied)}
}

type ProductionUnitReleased struct{ ev.BaseEvent }

func NewProductionUnitReleased(id vo.ID) ProductionUnitReleased {
	return ProductionUnitReleased{ev.NewBaseEvent(id, EventReleased)}
}

type ProductionUnitInPreparation struct{ ev.BaseEvent }

func NewProductionUnitInPreparation(id vo.ID) ProductionUnitInPreparation {
	return ProductionUnitInPreparation{ev.NewBaseEvent(id, EventInPreparation)}
}

type ProductionUnitArchived struct{ ev.BaseEvent }

func NewProductionUnitArchived(id vo.ID) ProductionUnitArchived {
	return ProductionUnitArchived{ev.NewBaseEvent(id, EventArchived)}
}

type ProductionUnitCreated struct {
	ev.BaseEvent
}

func NewProductionUnitCreated(
	id vo.ID,
) ProductionUnitCreated {
	return ProductionUnitCreated{
		BaseEvent: ev.NewBaseEvent(id, EventCreated),
	}
}

type ProductionUnitAttached struct {
	ev.BaseEvent
	ParentID vo.ID
}

func NewProductionUnitAttached(
	id vo.ID,
	parent vo.ID,
) ProductionUnitAttached {
	return ProductionUnitAttached{
		BaseEvent: ev.NewBaseEvent(id, EventAttached),
		ParentID:  parent,
	}
}
