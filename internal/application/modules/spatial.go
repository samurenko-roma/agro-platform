package modules

import (
	layoutsnapshotCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/spatial/layout_snapshot"
	productionunitCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/spatial/production_unit"
	layoutsnapshotQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/spatial/layout_snapshot"
	productionunitQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/spatial/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	spatialsnapshot "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/spatial/layout_snapshot"
	spatial "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/spatial/production_unit"
	"github.com/samurenkoroma/agro-platform/pkg/utils"
)

func MakeSpatialModule(uow uow.UnitOfWork, db uow.DB) Module {
	return Module{
		Commands: []*CommandCNF{
			{
				RouteName: "spatial.create_production_unit",
				Handler:   productionunitCmd.NewProductionUnitHandler(uow).Create,
				Decoder:   utils.DecodeJSON[productionunitCmd.CreateCommand],
			},
			{
				RouteName: "spatial.update_production_unit",
				Handler:   productionunitCmd.NewProductionUnitHandler(uow).Update,
				Decoder:   utils.DecodeJSON[productionunitCmd.UpdateCommand],
			},
			{
				RouteName: "spatial.configure_production_unit",
				Handler:   productionunitCmd.NewProductionUnitHandler(uow).Configure,
				Decoder:   utils.DecodeJSON[productionunitCmd.ConfigureCommand],
			},
			{
				RouteName: "spatial.archive_production_unit",
				Handler:   productionunitCmd.NewProductionUnitHandler(uow).Archive,
				Decoder:   utils.DecodeJSON[productionunitCmd.ArchiveCommand],
			},
			{
				RouteName: "spatial.create_layout_snapshot",
				Handler:   layoutsnapshotCmd.NewHandler(uow).Create,
				Decoder:   utils.DecodeJSON[layoutsnapshotCmd.CreateSnapshotCommand],
			},
			// spatial.move_production_unit  — TODO: следующая итерация
			// spatial.clone_production_unit — TODO: следующая итерация
		},
		Queries: []*QueryCNF{
			{
				RouteName: "spatial.get_production_unit",
				Handler:   productionunitQuery.NewGetOne(spatial.New(db)),
				Decoder:   utils.DecodeJSON[productionunitQuery.GetOneQuery],
			},
			{
				RouteName: "spatial.list_production_units",
				Handler:   productionunitQuery.NewListRoots(spatial.New(db)),
				Decoder:   utils.DecodeJSON[productionunitQuery.ListRootsQuery],
			},
			{
				RouteName: "spatial.get_production_unit_children",
				Handler:   productionunitQuery.NewListChildren(spatial.New(db)),
				Decoder:   utils.DecodeJSON[productionunitQuery.ListChildrenQuery],
			},
			{
				RouteName: "spatial.get_layout_snapshot",
				Handler:   layoutsnapshotQuery.NewGet(spatialsnapshot.New(db)),
				Decoder:   utils.DecodeJSON[layoutsnapshotQuery.GetQuery],
			},
			{
				RouteName: "spatial.get_latest_layout_snapshot",
				Handler:   layoutsnapshotQuery.NewGetLatest(spatialsnapshot.New(db)),
				Decoder:   utils.DecodeJSON[layoutsnapshotQuery.GetLatestQuery],
			},
		},
	}
}
