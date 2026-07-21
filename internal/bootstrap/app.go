package bootstrap

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/samurenkoroma/agro-platform/internal/application/commands"
	eventhandlers "github.com/samurenkoroma/agro-platform/internal/application/event_handlers"
	"github.com/samurenkoroma/agro-platform/internal/application/modules"
	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	unitOfWork "github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/jwt"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/messaging/inmemory"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/postgres"
	http2 "github.com/samurenkoroma/agro-platform/internal/interfaces/http"
	configs "github.com/samurenkoroma/agro-platform/internal/shared/config"
	"github.com/samurenkoroma/agro-platform/pkg/logger"
)

type App struct {
	DB            unitOfWork.DB
	CommandRouter commands.Router
	QueryRouter   queries.Router
	HTTPHandler   http.Handler
}

func Build(ctx context.Context, pool *pgxpool.Pool, conf *configs.Config) (*App, error) {
	log := logger.New(conf.Logger.Level, nil)
	ctx = logger.WithContext(ctx, log)

	log.Info("bootstrap: initializing application")

	bus := inmemory.NewInMemoryEventBus()
	uow := postgres.NewUnitOfWork(pool, bus)
	jwtService := jwt.NewService(conf.Auth)
	commandRouter := commands.NewRouter()
	queryRouter := queries.NewRouter()

	appModules := []modules.Module{
		modules.MakeAccountModule(uow, jwtService),
		modules.MakeAgronomyModule(uow, pool),
		modules.MakeSpatialModule(uow, pool),
		modules.MakeProductionModule(uow, pool),
		modules.MakeOperationsModule(uow, pool),
		modules.MakeInventoryModule(uow, pool),
		//modules.MakeWeatherModule(uow, pool),
	}

	for _, module := range appModules {
		module.RegisterQueries(queryRouter)
		module.RegisterCommands(commandRouter)
	}
	eventhandlers.RegisterAllocationHandlers(bus, uow)

	log.Info("bootstrap: modules registered", "count", len(appModules))

	httpHandler := http2.NewRouter(http2.RouterConfig{
		CommandRouter: commandRouter,
		QueryRouter:   queryRouter,
		Uow:           uow,
		JWTService:    jwtService,
		Logger:        log,
	})
	return &App{
		DB:            pool,
		CommandRouter: commandRouter,
		QueryRouter:   queryRouter,
		HTTPHandler:   httpHandler,
	}, nil
}
