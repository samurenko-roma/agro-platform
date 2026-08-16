package swaggerdocs

// docCreateCycle godoc
// @Summary Зарегистрировать цикл выращивания без немедленного размещения
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body growingcycleCmd.CreateCommand true "Культура, сорт, протокол, метод"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/production.create_cycle [post]
func docCreateCycle() {}

// docStartCycle godoc
// @Summary Создать цикл выращивания И сразу разместить его на физических узлах (+опционально посевы)
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body growingcycleCmd.StartGrowingCycleCMD true "Культура, сорт, метод, обязательные allocations, опциональные plantings"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/production.start_cycle [post]
func docStartCycle() {}

// docAllocateProductionUnit godoc
// @Summary Разместить (доп.) цикл выращивания на физическом узле
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body allocationCmd.AllocateProductionUnitCommand true "ID цикла, ID узла, площадь"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 404 {object} response.CommandResponse[any] "цикл не найден в вашей организации"
// @Router /api/commands/production.allocate_production_unit [post]
func docAllocateProductionUnit() {}

// docChangeAllocation godoc
// @Summary Изменить размещение (площадь/даты, либо перенести на другой узел)
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body allocationCmd.ChangeAllocationCommand true "ID размещения, новые параметры"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/production.change_allocation [post]
func docChangeAllocation() {}

// docReleaseAllocation godoc
// @Summary Завершить размещение (освободить физический узел)
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body allocationCmd.ReleaseAllocationCommand true "ID размещения, дата освобождения (опционально — по умолчанию сейчас)"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 409 {object} response.CommandResponse[any] "уже освобождено"
// @Router /api/commands/production.release_allocation [post]
func docReleaseAllocation() {}

// docPlantingRegister godoc
// @Summary Зарегистрировать акт посева/посадки
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body planting.RegisterPlantingCommand true "ID цикла, дата, количество"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 404 {object} response.CommandResponse[any] "цикл не найден в вашей организации"
// @Router /api/commands/production.planting_register [post]
func docPlantingRegister() {}

// docPlantingChange godoc
// @Summary Изменить ранее зарегистрированный посев
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body planting.ChangePlantingCommand true "ID посева, новые дата/количество"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/production.planting_change [post]
func docPlantingChange() {}

// docHarvestRegister godoc
// @Summary Зарегистрировать сбор урожая
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body harvest.RegisterHarvestCommand true "ID цикла, дата, количество"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 404 {object} response.CommandResponse[any] "цикл не найден в вашей организации"
// @Router /api/commands/production.harvest_register [post]
func docHarvestRegister() {}

// docHarvestChange godoc
// @Summary Изменить ранее зарегистрированный сбор урожая
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body harvest.ChangeHarvestCommand true "ID харвеста, новые дата/количество"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/production.harvest_change [post]
func docHarvestChange() {}

// docGetGrowingCycle godoc
// @Summary Получить один цикл выращивания (сводка: площадь, посажено, собрано)
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body growingcycleQuery.GetOneQuery true "ID цикла"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/production.get_growing_cycle [post]
func docGetGrowingCycle() {}

// docListGrowingCycles godoc
// @Summary Список циклов выращивания, сгруппированных по культуре, с аллокациями и прогрессом
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/production.list_growing_cycles [post]
func docListGrowingCycles() {}

// docSummaryGrowingCycles godoc
// @Summary Сводка по одному циклу: занятая площадь, посажено, собрано
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body growingcycleQuery.SummaryQuery true "ID цикла"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/production.summary_growing_cycles [post]
func docSummaryGrowingCycles() {}

// docGetAllocation godoc
// @Summary Получить размещение по ID (с кодом физического узла)
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body allocationQuery.GetOneQuery true "ID размещения"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/production.get_allocation [post]
func docGetAllocation() {}

// docListAllocations godoc
// @Summary Список всех размещений организации
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/production.list_allocations [post]
func docListAllocations() {}

// docHelpers godoc
// @Summary Справочник допустимых статусов/стадий/методов цикла выращивания (для выпадающих списков UI)
// @Tags production
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/production.helpers [post]
func docHelpers() {}
