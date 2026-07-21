package providers

import (
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	weatherrepo "github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
	pgweather "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/postgres/weather"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type weatherProvider struct {
	db        uow.DB
	locations weatherrepo.WeatherLocationRepository
	records   weatherrepo.WeatherRecordRepository
}

func NewWeatherProvider(db uow.DB) repository.RepositoryProvider {
	return &weatherProvider{db: db}
}

func (p *weatherProvider) ProviderName() string { return "weather" }

func (p *weatherProvider) Locations() weatherrepo.WeatherLocationRepository {
	if p.locations == nil {
		p.locations = pgweather.NewLocationRepository(p.db)
	}
	return p.locations
}

func (p *weatherProvider) Records() weatherrepo.WeatherRecordRepository {
	if p.records == nil {
		p.records = pgweather.NewRecordRepository(p.db)
	}
	return p.records
}

var _ weatherrepo.WeatherProvider = (*weatherProvider)(nil)
