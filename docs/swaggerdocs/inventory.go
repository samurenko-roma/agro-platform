package swaggerdocs

import (
	itemCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/inventory/item"
	warehouseCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/inventory/warehouse"
	itemQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/inventory/item"
	movQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/inventory/movement"
	warehouseQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/inventory/warehouse"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

var (
	_ = itemCmd.CreateItemCommand{}
	_ = warehouseCmd.CreateWarehouseCommand{}
	_ = warehouseQuery.ListQuery{}
	_ = movQuery.ListQuery{}
	_ = itemQuery.ListQuery{}
	_ = response.CommandResponse{}
)

// docCreateItem godoc
// @Summary Создать позицию склада (семена, удобрения, субстрат и т.д.)
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.create_item", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.CreateItemCommand true "Название, тип, единица измерения, склад (опционально)"
// @Success 200 {object} response.CommandResponse
// @Router /api/commands/inventory.create_item [post]
func docCreateItem() {}

// docReceive godoc
// @Summary Оприходовать поступление на склад
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.receive", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.ReceiveCommand true "ID позиции, количество, примечание"
// @Success 200 {object} response.CommandResponse
// @Router /api/commands/inventory.receive [post]
func docReceive() {}

// docReserve godoc
// @Summary Зарезервировать позицию (например, под цикл выращивания)
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.reserve", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.ReserveCommand true "ID позиции, количество, ссылка (тип+ID)"
// @Success 200 {object} response.CommandResponse
// @Failure 409 {object} response.CommandResponse "недостаточно свободного остатка"
// @Router /api/commands/inventory.reserve [post]
func docReserve() {}

// docConsume godoc
// @Summary Списать зарезервированное (реальный расход)
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.consume", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.ConsumeCommand true "ID позиции, количество, ссылка (тип+ID)"
// @Success 200 {object} response.CommandResponse
// @Failure 409 {object} response.CommandResponse "недостаточно зарезервированного остатка"
// @Router /api/commands/inventory.consume [post]
func docConsume() {}

// docMarkLost godoc
// @Summary Списать позицию как утраченную (порча, просыпано и т.д.)
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.mark_lost", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.MarkLostCommand true "ID позиции, количество, примечание"
// @Success 200 {object} response.CommandResponse
// @Router /api/commands/inventory.mark_lost [post]
func docMarkLost() {}

// docTransfer godoc
// @Summary Переместить позицию между складами
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.transfer", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.TransferCommand true "ID позиции, количество, ID складов отправитель/получатель"
// @Success 200 {object} response.CommandResponse
// @Router /api/commands/inventory.transfer [post]
func docTransfer() {}

// docCreateWarehouse godoc
// @Summary Создать склад
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.create_warehouse", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body warehouseCmd.CreateWarehouseCommand true "Название, код (опционально)"
// @Success 200 {object} response.CommandResponse
// @Router /api/commands/inventory.create_warehouse [post]
func docCreateWarehouse() {}

// docArchiveWarehouse godoc
// @Summary Архивировать склад
// @Description Реальный вызов: POST /api/commands, {"command": "inventory.archive_warehouse", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body warehouseCmd.ArchiveWarehouseCommand true "ID склада"
// @Success 200 {object} response.CommandResponse
// @Router /api/commands/inventory.archive_warehouse [post]
func docArchiveWarehouse() {}

// docGetItem godoc
// @Summary Получить позицию склада по ID
// @Description Реальный вызов: POST /api/queries, {"query": "inventory.get_item", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemQuery.GetOneQuery true "ID позиции"
// @Success 200 {object} response.CommandResponse
// @Router /api/queries/inventory.get_item [post]
func docGetItem() {}

// docListItems godoc
// @Summary Список позиций склада (опционально — по конкретному складу)
// @Description Реальный вызов: POST /api/queries, {"query": "inventory.list_items", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemQuery.ListQuery false "warehouseId — опционально"
// @Success 200 {object} response.CommandResponse
// @Router /api/queries/inventory.list_items [post]
func docListItems() {}

// docListMovements godoc
// @Summary История движений (опционально — по конкретной позиции)
// @Description Реальный вызов: POST /api/queries, {"query": "inventory.list_movements", "data": <тело ниже>}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body movQuery.ListQuery false "itemId — опционально"
// @Success 200 {object} response.CommandResponse
// @Router /api/queries/inventory.list_movements [post]
func docListMovements() {}

// docListWarehouses godoc
// @Summary Список складов организации
// @Description Реальный вызов: POST /api/queries, {"query": "inventory.list_warehouses", "data": {}}
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse
// @Router /api/queries/inventory.list_warehouses [post]
func docListWarehouses() {}
