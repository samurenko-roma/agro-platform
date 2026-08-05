package swaggerdocs

import (
	opCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/operations/operation_event"
	taskCmd "github.com/samurenkoroma/agro-platform/internal/application/commands/operations/task"
	opQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/operations/operation_event"
	taskQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/operations/task"
	tlQuery "github.com/samurenkoroma/agro-platform/internal/application/queries/operations/timeline"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

var (
	_ = taskCmd.CreateTaskCommand{}
	_ = tlQuery.GetQuery{}
	_ = taskQuery.GetOneQuery{}
	_ = opCmd.RecordOperationCommand{}
	_ = opQuery.ListQuery{}
	_ = response.CommandResponse[any]{}
)

// docCreateTask godoc
// @Summary Создать задачу (наряд на выполнение работ)
// @Description Реальный вызов: POST /api/commands, {"command": "operations.create_task", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/commands, {"command": "operations.assign_task", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/commands, {"command": "operations.start_task", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/commands, {"command": "operations.complete_task", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/commands, {"command": "operations.cancel_task", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/commands, {"command": "operations.record_operation", "data": <тело ниже>}. Автоматически добавляется в таймлайн цикла/фермы.
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
// @Description Реальный вызов: POST /api/queries, {"query": "operations.get_task", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/queries, {"query": "operations.list_tasks", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/queries, {"query": "operations.get_timeline", "data": <тело ниже>}
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
// @Description Реальный вызов: POST /api/queries, {"query": "operations.list_operations", "data": <тело ниже>}
// @Tags operations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body opQuery.ListQuery false "growingCycleId — опционально"
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/operations.list_operations [post]
func docListOperations() {}
