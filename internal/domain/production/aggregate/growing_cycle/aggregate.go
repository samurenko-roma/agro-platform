package growingcycle

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type GrowingCycle struct {
	ev.BaseAggregate
	ID vo.ID

	FarmID vo.ID

	CropID    vo.ID
	VarietyID *vo.ID

	ProtocolID *vo.ID

	Name string
	Code string

	Method ProductionMethod

	Status CycleStatus
	Stage  CycleStage

	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(orgId, cropID vo.ID, varietyId, protocolId *vo.ID, name, code string, method ProductionMethod) *GrowingCycle {
	now := time.Now()
	root := &GrowingCycle{
		ID:         vo.NewID(),
		FarmID:     orgId,
		CropID:     cropID,
		VarietyID:  varietyId,
		ProtocolID: protocolId,
		Name:       name,
		Code:       code,
		Method:     method,
		Status:     StatusPlanned,
		Stage:      StagePlanning,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	root.AddEvent(NewCycleCreated(root.ID))

	return root
}

func (gc *GrowingCycle) ChangeState(state CycleStage) {
	if state == "" {
		gc.Stage = StagePlanning
		return
	}
	gc.Stage = state
	gc.UpdatedAt = time.Now()
}

func (gc *GrowingCycle) ChangeStatus(status CycleStatus) {
	if status == "" {
		gc.Status = StatusPlanned
		return
	}
	gc.Status = status
	gc.UpdatedAt = time.Now()
}

// Start — единственный доменный способ активировать цикл. В отличие от
// ChangeStatus, не принимает произвольную строку от клиента: "запустить"
// всегда означает переход в StatusActive, и эмитит CycleStarted.
func (gc *GrowingCycle) Start() {
	gc.Status = StatusActive
	gc.UpdatedAt = time.Now()
	gc.AddEvent(NewCycleStarted(gc.ID))
}
