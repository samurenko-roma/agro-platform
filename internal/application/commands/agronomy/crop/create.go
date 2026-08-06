package crop

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/agronomy/aggregate/crop"
	agronomy "github.com/samurenkoroma/agro-platform/internal/domain/agronomy/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type CreateCropCommand struct {
	Name     string `json:"name" validate:"required"`
	Category string `json:"category" required:"true"`
	Family   string `json:"family" required:"true"`

	Agronomy    crop.AgronomyProfile `json:"agronomy"`
	Description *string              `json:"description,omitempty"`
}

func (h *Handler) Create(ctx context.Context, payload any) (any, error) {
	cmd := payload.(*CreateCropCommand)

	return h.uow.Execute(ctx, providers.NewAgronomyProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {
		agronomyProvider, ok := provider.(agronomy.AgronomyProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		exist, _ := agronomyProvider.Crops().Exists(ctx, cmd.Name)
		if exist {
			return nil, ErrCropAlreadyExist
		}

		root := crop.New(cmd.Name, crop.CropCategory(cmd.Category), cmd.Family, cmd.Agronomy)

		if cmd.Description != nil {
			root.Metadata["description"] = cmd.Description
		}
		if err := agronomyProvider.Crops().Save(ctx, root); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(root)

		return response.Id(root.ID), nil
	})
}
