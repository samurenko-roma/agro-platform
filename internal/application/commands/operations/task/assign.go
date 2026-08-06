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

type AssignTaskCommand struct {
	TaskID string `json:"taskId" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

func (h *Handler) Assign(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*AssignTaskCommand)
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
		t.Assign(vo.ID(cmd.UserID))
		if err := ops.Tasks().Save(ctx, t); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(t)
		return response.Id(t.ID), nil
	})
}
