package growingcycle

import (
	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/event"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

const (
	EventCreated = "growing_cycle.created"

	EventStarted = "growing_cycle.started"

	EventPaused = "growing_cycle.paused"

	EventResumed = "growing_cycle.resumed"

	EventFailed = "growing_cycle.failed"

	EventArchived       = "growing_cycle.archived"
	EventHarvestStarted = "growing_cycle.harvest.started"

	EventPartialHarvest = "growing_cycle.harvest.partial"

	EventHarvestCompleted = "growing_cycle.harvest.completed"
)

type CycleCreated struct {
	ev.BaseEvent
}

func NewCycleCreated(id vo.ID) CycleCreated {
	return CycleCreated{BaseEvent: ev.NewBaseEvent(id, EventCreated)}
}

type CycleStarted struct{ ev.BaseEvent }
type CyclePaused struct{ ev.BaseEvent }
type CycleResumed struct{ ev.BaseEvent }
type CycleFailed struct{ ev.BaseEvent }
type CycleArchived struct{ ev.BaseEvent }

type CycleHarvestStarted struct{ ev.BaseEvent }
type CyclePartialHarvest struct {
	RecordId vo.ID
	ev.BaseEvent
}
type CycleHarvestCompleted struct{ ev.BaseEvent }

func NewHarvestStarted(id vo.ID) CycleHarvestStarted {
	return CycleHarvestStarted{BaseEvent: ev.NewBaseEvent(id, EventHarvestStarted)}
}

func NewPartialHarvest(id vo.ID, recordId vo.ID) CyclePartialHarvest {
	return CyclePartialHarvest{RecordId: recordId, BaseEvent: ev.NewBaseEvent(id, EventPartialHarvest)}
}

func NewHarvestCompleted(id vo.ID) CycleHarvestCompleted {
	return CycleHarvestCompleted{BaseEvent: ev.NewBaseEvent(id, EventHarvestCompleted)}
}

func NewCycleStarted(id vo.ID) CycleStarted {
	return CycleStarted{BaseEvent: ev.NewBaseEvent(id, EventStarted)}
}

func NewCyclePaused(id vo.ID) CyclePaused {
	return CyclePaused{BaseEvent: ev.NewBaseEvent(id, EventPaused)}
}

func NewCycleResumed(id vo.ID) CycleResumed {
	return CycleResumed{BaseEvent: ev.NewBaseEvent(id, EventResumed)}
}

func NewCycleFailed(id vo.ID) CycleFailed {
	return CycleFailed{BaseEvent: ev.NewBaseEvent(id, EventFailed)}
}

func NewCycleArchived(id vo.ID) CycleArchived {
	return CycleArchived{BaseEvent: ev.NewBaseEvent(id, EventArchived)}
}
