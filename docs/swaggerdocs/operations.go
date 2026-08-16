package swaggerdocs

// docCreateTask godoc
// @Summary Создать задачу (наряд на выполнение работ)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskCmd.CreateTaskCommand true "Заголовок, приоритет, привязка к циклу/узлу (опционально)"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/operations.create_task [post]
func docCreateTask() {}

// docAssignTask godoc
// @Summary Назначить задачу исполнителю
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskCmd.AssignTaskCommand true "ID задачи, ID пользователя"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/operations.assign_task [post]
func docAssignTask() {}

// docStartTask godoc
// @Summary Взять задачу в работу (TODO → IN_PROGRESS)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskCmd.TaskIDCommand true "ID задачи"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/operations.start_task [post]
func docStartTask() {}

// docCompleteTask godoc
// @Summary Завершить задачу (→ DONE)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskCmd.TaskIDCommand true "ID задачи"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/operations.complete_task [post]
func docCompleteTask() {}

// docCancelTask godoc
// @Summary Отменить задачу (→ CANCELLED)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskCmd.TaskIDCommand true "ID задачи"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/operations.cancel_task [post]
func docCancelTask() {}

// docRecordOperation godoc
// @Summary Зафиксировать агрономическую операцию (полив, удобрение, обработка и т.д.)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body opCmd.RecordOperationCommand true "Тип операции, привязка к узлу/циклу, произвольный payload"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/commands/operations.record_operation [post]
func docRecordOperation() {}

// docGetTask godoc
// @Summary Получить задачу по ID
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskQuery.GetOneQuery true "ID задачи"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/operations.get_task [post]
func docGetTask() {}

// docListTasks godoc
// @Summary Список задач фермы (опционально — по циклу выращивания)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskQuery.ListQuery false "growingCycleId — опционально"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/operations.list_tasks [post]
func docListTasks() {}

// docGetTimeline godoc
// @Summary Таймлайн операций/задач (по циклу или общий по ферме)
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tlQuery.GetQuery false "growingCycleId — опционально, иначе общий таймлайн фермы"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/operations.get_timeline [post]
func docGetTimeline() {}

// docListOperations godoc
// @Summary Список зафиксированных агрономических операций
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body opQuery.ListQuery false "growingCycleId — опционально"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/operations.list_operations [post]
func docListOperations() {}
