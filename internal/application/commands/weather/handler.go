package weather

import (
	weatherservice "github.com/samurenkoroma/agro-platform/internal/application/services/weather"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
)

type Handler struct {
	uow     uow.UnitOfWork
	service *weatherservice.WeatherService
}

func NewHandler(uow uow.UnitOfWork, svc *weatherservice.WeatherService) *Handler {
	return &Handler{uow: uow, service: svc}
}
