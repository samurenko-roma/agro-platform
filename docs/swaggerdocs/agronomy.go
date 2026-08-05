package swaggerdocs

import (
	"github.com/samurenkoroma/agro-platform/internal/application/commands/agronomy/crop"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/agronomy/season"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/agronomy/variety"
	getcrop "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/crop/get_crop"
	listcrops "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/crop/list_crops"
	getvariety "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/variety/get_variety"
	listvarieties "github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/variety/list_varieties"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

var (
	_ = crop.CreateCropCommand{}
	_ = variety.CreateVarietyCommand{}
	_ = season.CreateSeasonCmd{}
	_ = getcrop.Query{}
	_ = listcrops.Query{}
	_ = getvariety.Query{}
	_ = listvarieties.Query{}

	_ = response.CommandResponse{}
)

// docCreateCrop godoc
// @Summary Создать культуру в справочнике (Crop) — биология вида: температуры, GDD, водопотребление
// @Description Реальный вызов: POST /api/commands, {"command": "agronomy.create_crop", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body crop.CreateCropCommand true "Название, категория, семейство, агрономический профиль"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 409 {object} response.CommandResponse[any] "культура с таким именем уже существует"
// @Router /api/commands/agronomy.create_crop [post]
func docCreateCrop() {}

// docCreateVariety godoc
// @Summary Создать сорт культуры (Variety) — конкретный коммерческий сорт с зрелостью/спейсингом
// @Description Реальный вызов: POST /api/commands, {"command": "agronomy.create_variety", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body variety.CreateVarietyCommand true "Название, ID культуры, дни до созревания"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 404 {object} response.CommandResponse[any] "культура не найдена"
// @Failure 409 {object} response.CommandResponse[any] "сорт с таким именем уже существует у этой культуры"
// @Router /api/commands/agronomy.create_variety [post]
func docCreateVariety() {}

// docCreateSeason godoc
// @Summary Создать агрономический сезон (период планирования)
// @Description Реальный вызов: POST /api/commands, {"command": "agronomy.create_season", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body season.CreateSeasonCmd true "Название, даты начала/конца, статус"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/agronomy.create_season [post]
func docCreateSeason() {}

// docGetCrop godoc
// @Summary Получить культуру по ID
// @Description Реальный вызов: POST /api/queries, {"query": "agronomy.get_crop", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body getcrop.Query true "ID культуры"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/agronomy.get_crop [post]
func docGetCrop() {}

// docListCrops godoc
// @Summary Список культур (с фильтром по категории)
// @Description Реальный вызов: POST /api/queries, {"query": "agronomy.list_crops", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body listcrops.Query false "Фильтр по категориям/поиску (опционально)"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/agronomy.list_crops [post]
func docListCrops() {}

// docGetVariety godoc
// @Summary Получить сорт по ID
// @Description Реальный вызов: POST /api/queries, {"query": "agronomy.get_variety", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body getvariety.Query true "ID сорта"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/agronomy.get_variety [post]
func docGetVariety() {}

// docListVarieties godoc
// @Summary Список сортов культуры
// @Description Реальный вызов: POST /api/queries, {"query": "agronomy.list_varieties", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body listvarieties.Query true "ID культуры"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/agronomy.list_varieties [post]
func docListVarieties() {}

// docListSeasons godoc
// @Summary Список сезонов организации
// @Description Реальный вызов: POST /api/queries, {"query": "agronomy.list_seasons", "data": {}}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/agronomy.list_seasons [post]
func docListSeasons() {}
