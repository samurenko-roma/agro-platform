package task

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	opsrepo "github.com/samurenkoroma/agro-platform/internal/domain/operations/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type TaskIDCommand struct {
	TaskID string `json:"taskId" validate:"required"`
}

type action string

const (
	actionStart    action = "start"
	actionComplete action = "complete"
	actionCancel   action = "cancel"
)

func (h *Handler) Start(ctx context.Context, payload any) (any, error) {
	return h.changeStatus(ctx, payload, actionStart)
}

func (h *Handler) Complete(ctx context.Context, payload any) (any, error) {
	return h.changeStatus(ctx, payload, actionComplete)
}

func (h *Handler) Cancel(ctx context.Context, payload any) (any, error) {
	return h.changeStatus(ctx, payload, actionCancel)
}

func (h *Handler) changeStatus(ctx context.Context, payload any, a action) (any, error) {
	cmd, ok := payload.(*TaskIDCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	if _, ok := ctx.Value("organization_id").(string); !ok {
		return nil, errors.New("organization_id is required")
	}
	return h.uow.Execute(ctx, providers.NewOperationsProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		ops, ok := p.(opsrepo.OperationsProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}
		t, err := ops.Tasks().GetByID(ctx, vo.ID(cmd.TaskID))
		if err != nil {
			return nil, err
		}
		switch a {
		case actionStart:
			t.Start()
		case actionComplete:
			t.Complete()
		case actionCancel:
			t.Cancel()
		}
		if err := ops.Tasks().Save(ctx, t); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(t)
		return response.Id(t.ID), nil
	})
}
