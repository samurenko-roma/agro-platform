package task

import (
	"context"
	"errors"
	"time"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	taskdomain "github.com/samurenkoroma/agro-platform/internal/domain/operations/aggregate/task"
	opsrepo "github.com/samurenkoroma/agro-platform/internal/domain/operations/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type CreateTaskCommand struct {
	Title            string     `json:"title" validate:"required"`
	Description      *string    `json:"description"`
	OperationType    *string    `json:"operationType"`
	ProductionUnitID *string    `json:"productionUnitId"`
	GrowingCycleID   *string    `json:"growingCycleId"`
	Priority         string     `json:"priority"`
	DueDate          *time.Time `json:"dueDate"`
}

func (h *Handler) Create(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*CreateTaskCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}

	return h.uow.Execute(ctx, providers.NewOperationsProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		ops, ok := p.(opsrepo.OperationsProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		t := taskdomain.New(vo.ID(orgID), cmd.Title)
		t.Description = cmd.Description
		t.DueDate = cmd.DueDate

		if cmd.Priority != "" {
			t.Priority = taskdomain.Priority(cmd.Priority)
		}
		if cmd.ProductionUnitID != nil {
			id := vo.ID(*cmd.ProductionUnitID)
			t.ProductionUnitID = &id
		}
		if cmd.GrowingCycleID != nil {
			id := vo.ID(*cmd.GrowingCycleID)
			t.GrowingCycleID = &id
		}

		if err := ops.Tasks().Save(ctx, t); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(t)
		return response.Id(t.ID), nil
	})
}
