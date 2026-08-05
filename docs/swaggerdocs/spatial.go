package swaggerdocs

import (
	productionunitCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/spatial/production_unit"
	productionunitQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/spatial/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

var (
	_ = productionunitCmd.CreateCommand{}
	_ = productionunitCmd.UpdateCommand{}
	_ = productionunitCmd.ConfigureCommand{}
	_ = productionunitCmd.ArchiveCommand{}
	_ = productionunitQuery.DTO{}

	_ = productionunitQuery.GetOneQuery{}

	_ = response.CommandResponse[any, any]{}
)

// docCreateProductionUnit godoc
// @Summary Создать узел (поле, участок, грядку, теплицу и т.д.)
// @Description Реальный вызов: POST /api/commands, {"command": "spatial.create_production_unit", "data": <тело ниже>}
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productionunitCmd.CreateCommand true "Параметры узла"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 400 {object} response.CommandResponse[any]
// @Router /api/commands/spatial.create_production_unit [post]
func docCreateProductionUnit() {}

// docUpdateProductionUnit godoc
// @Summary Обновить JSON-схему узла
// @Description Реальный вызов: POST /api/commands, {"command": "spatial.update_production_unit", "data": <тело ниже>}
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productionunitCmd.UpdateCommand true "ID узла и схема"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/spatial.update_production_unit [post]
func docUpdateProductionUnit() {}

// docConfigureProductionUnit godoc
// @Summary Сгенерировать дочерние узлы по схеме (например, грядки в теплице)
// @Description Реальный вызов: POST /api/commands, {"command": "spatial.configure_production_unit", "data": <тело ниже>}
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productionunitCmd.ConfigureCommand true "ID родителя и схема раскладки"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/spatial.configure_production_unit [post]
func docConfigureProductionUnit() {}

// docArchiveProductionUnit godoc
// @Summary Архивировать узел
// @Description Реальный вызов: POST /api/commands, {"command": "spatial.archive_production_unit", "data": <тело ниже>}
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productionunitCmd.ArchiveCommand true "ID узла"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 409 {object} response.CommandResponse[any] "есть неархивные дети"
// @Router /api/commands/spatial.archive_production_unit [post]
func docArchiveProductionUnit() {}

// docGetProductionUnit godoc
// @Summary Получить узел вместе с поддеревом
// @Description Реальный вызов: POST /api/queries, {"query": "spatial.get_production_unit", "data": <тело ниже>}
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productionunitQuery.GetOneQuery true "ID узла"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/spatial.get_production_unit [post]
func docGetProductionUnit() {}

// docListProductionUnits godoc
// @Summary Список корневых узлов организации (с вложенными деревьями)
// @Description Реальный вызов: POST /api/queries, {"query": "spatial.list_production_units", "data": {}}
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[[]productionunitQuery.DTO, nil]
// @Router /api/queries/spatial.list_production_units [post]
func docListProductionUnits() {}
