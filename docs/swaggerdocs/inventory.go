package swaggerdocs

// docCreateItem godoc
// @Summary Создать позицию склада (семена, удобрения, субстрат и т.д.)
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.CreateItemCommand true "Название, тип, единица измерения, склад (опционально)"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/inventory.create_item [post]
func docCreateItem() {}

// docReceive godoc
// @Summary Оприходовать поступление на склад
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.ReceiveCommand true "ID позиции, количество, примечание"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/inventory.receive [post]
func docReceive() {}

// docReserve godoc
// @Summary Зарезервировать позицию (например, под цикл выращивания)
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.ReserveCommand true "ID позиции, количество, ссылка (тип+ID)"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 409 {object} response.CommandResponse[any] "недостаточно свободного остатка"
// @Router /api/commands/inventory.reserve [post]
func docReserve() {}

// docConsume godoc
// @Summary Списать зарезервированное (реальный расход)
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.ConsumeCommand true "ID позиции, количество, ссылка (тип+ID)"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 409 {object} response.CommandResponse[any] "недостаточно зарезервированного остатка"
// @Router /api/commands/inventory.consume [post]
func docConsume() {}

// docMarkLost godoc
// @Summary Списать позицию как утраченную (порча, просыпано и т.д.)
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.MarkLostCommand true "ID позиции, количество, примечание"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/inventory.mark_lost [post]
func docMarkLost() {}

// docTransfer godoc
// @Summary Переместить позицию между складами
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemCmd.TransferCommand true "ID позиции, количество, ID складов отправитель/получатель"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/inventory.transfer [post]
func docTransfer() {}

// docCreateWarehouse godoc
// @Summary Создать склад
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body warehouseCmd.CreateWarehouseCommand true "Название, код (опционально)"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/inventory.create_warehouse [post]
func docCreateWarehouse() {}

// docArchiveWarehouse godoc
// @Summary Архивировать склад
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body warehouseCmd.ArchiveWarehouseCommand true "ID склада"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/inventory.archive_warehouse [post]
func docArchiveWarehouse() {}

// docGetItem godoc
// @Summary Получить позицию склада по ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemQuery.GetOneQuery true "ID позиции"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/inventory.get_item [post]
func docGetItem() {}

// docListItems godoc
// @Summary Список позиций склада (опционально — по конкретному складу)
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body itemQuery.ListQuery false "warehouseId — опционально"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/inventory.list_items [post]
func docListItems() {}

// docListMovements godoc
// @Summary История движений (опционально — по конкретной позиции)
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body movQuery.ListQuery false "itemId — опционально"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/inventory.list_movements [post]
func docListMovements() {}

// docListWarehouses godoc
// @Summary Список складов организации
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/inventory.list_warehouses [post]
func docListWarehouses() {}
