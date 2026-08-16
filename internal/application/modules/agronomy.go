package modules

import (
	cropCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/agronomy/crop"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/agronomy/season"
	varietyCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/agronomy/variety"
	cropQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/crop"
	"github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/season/list_seasons"
	listvarieties "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/variety"
	varietyQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/variety"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	cropProjection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/agronomy/crop"
	varietyProjection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/agronomy/variety"
	"github.com/samurenkoroma/agro-platform/pkg/utils"
)

func MakeAgronomyModule(uow uow.UnitOfWork, db uow.DB) Module {
	return Module{
		Commands: []*CommandCNF{
			{
				RouteName: "agronomy.create_crop",
				Handler:   cropCmd.NewHandler(uow).Create,
				Decoder:   utils.DecodeJSON[cropCmd.CreateCropCommand],
			},
			{
				RouteName: "agronomy.create_variety",
				Handler:   varietyCmd.NewHandler(uow).Create,
				Decoder:   utils.DecodeJSON[varietyCmd.CreateVarietyCommand],
			},
			{
				RouteName: "agronomy.create_season",
				Handler:   season.NewHandler(uow).Create,
				Decoder:   utils.DecodeJSON[season.CreateSeasonCmd],
			},
		},
		Queries: []*QueryCNF{
			{
				RouteName: "agronomy.get_crop",
				Handler:   cropQuery.NewGetOneHandler(cropProjection.New(db)),
				Decoder:   utils.DecodeJSON[cropQuery.OneQuery],
			},
			{
				RouteName: "agronomy.list_crops",
				Handler:   cropQuery.NewGetListHandler(cropProjection.New(db)),
				Decoder:   utils.DecodeJSON[cropQuery.ListQuery],
			},
			{
				RouteName: "agronomy.get_variety",
				Handler:   listvarieties.NewGetOneHandler(varietyProjection.New(db)),
				Decoder:   utils.DecodeJSON[varietyQuery.OneQuery],
			},
			{
				RouteName: "agronomy.list_varieties",
				Handler:   listvarieties.New(varietyProjection.New(db)),
				Decoder:   utils.DecodeJSON[varietyQuery.ListQuery],
			},
			{
				RouteName: "agronomy.list_seasons",
				Handler:   listseasons.New(db),
				Decoder:   utils.DecodeJSON[listseasons.Query],
			},
		},
	}
}
