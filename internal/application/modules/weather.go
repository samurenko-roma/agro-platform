package modules

//
//import (
//	weatherCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/weather"
//	weatherQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/weather"
//	weatherservice "github.com/samurenkoroma/agro-platform/internal/application/services/weather"
//	"github.com/samurenkoroma/agro-platform/internal/application/uow"
//	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
//	pgweather "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/weather"
//	weatherprojection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/weather"
//	"github.com/samurenkoroma/agro-platform/internal/infrastructure/weather/openmeteo"
//	"github.com/samurenkoroma/agro-platform/pkg/utils"
//)
//
//func MakeWeatherModule(u uow.UnitOfWork, db uow.DB) Module {
//	// Репозитории для WeatherService (без транзакции — сервис читает/пишет напрямую)
//	locations := pgweather.NewLocationRepository(db)
//	records := pgweather.NewRecordRepository(db)
//
//	// Провайдеры данных
//	primary := openmeteo.NewProvider()
//	// Sensor provider добавляется через WithFallback если нужен
//	_ = providers.NewWeatherProvider
//
//	svc := weatherservice.New(primary, nil, records, locations)
//	h := weatherCmd.NewHandler(u, svc)
//	qh := weatherQuery.NewQueryHandler(weatherprojection.New(db))
//
//	return Module{
//		Commands: []*CommandCNF{
//			{
//				RouteName: "weather.create_location",
//				Handler:   h.CreateLocation,
//				Decoder:   utils.DecodeJSON[weatherCmd.CreateLocationCommand],
//			},
//			{
//				RouteName: "weather.sync_current",
//				Handler:   h.SyncCurrent,
//				Decoder:   utils.DecodeJSON[weatherCmd.SyncCurrentCommand],
//			},
//			{
//				RouteName: "weather.sync_forecast",
//				Handler:   h.SyncForecast,
//				Decoder:   utils.DecodeJSON[weatherCmd.SyncForecastCommand],
//			},
//			{
//				RouteName: "weather.sync_historical",
//				Handler:   h.SyncHistorical,
//				Decoder:   utils.DecodeJSON[weatherCmd.SyncHistoricalCommand],
//			},
//		},
//		Queries: []*QueryCNF{
//			{
//				RouteName: "weather.get_current",
//				Handler:   qh.GetCurrent,
//				Decoder:   utils.DecodeJSON[weatherQuery.GetCurrentQuery],
//			},
//			{
//				RouteName: "weather.list_forecast",
//				Handler:   qh.ListForecast,
//				Decoder:   utils.DecodeJSON[weatherQuery.ListForecastQuery],
//			},
//			{
//				RouteName: "weather.list_historical",
//				Handler:   qh.ListHistorical,
//				Decoder:   utils.DecodeJSON[weatherQuery.ListHistoricalQuery],
//			},
//			{
//				RouteName: "weather.list_locations",
//				Handler:   qh.ListLocations,
//				Decoder:   utils.DecodeJSON[weatherQuery.ListLocationsQuery],
//			},
//		},
//	}
//}
