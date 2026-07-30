package modules

import (
	productionunitCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/spatial/production_unit"
	productionunitQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/spatial/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
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
			//{
			//	RouteName: "spatial.archive_production_unit",
			//	//Handler:   createproductionunit.New(uow).Create,
			//	//Decoder:   utils.DecodeJSON[createproductionunit.CreateCommand],
			//},
			//{
			//	RouteName: "spatial.move_production_unit",
			//	//Handler:   createproductionunit.New(uow).Create,
			//	//Decoder:   utils.DecodeJSON[createproductionunit.CreateCommand],
			//},
			//{
			//	RouteName: "spatial.clone_production_unit",
			//	//Handler:   createproductionunit.New(uow).Create,
			//	//Decoder:   utils.DecodeJSON[createproductionunit.CreateCommand],
			//},
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
			//{
			//	RouteName: "spatial.get_production_unit_tree",
			//},
			//{
			//	RouteName: "spatial.get_production_unit_children",
			//},
		},
	}
}
