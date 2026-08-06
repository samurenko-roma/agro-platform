package task_test

import (
	"testing"

	task2 "github.com/samurenkoroma/agro-platform/internal/application/commands/operations/task"
	"github.com/samurenkoroma/agro-platform/internal/domain/operations/aggregate/task"
	opsrepo "github.com/samurenkoroma/agro-platform/internal/domain/operations/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/testutil"
)

func createAndGet(t *testing.T) (*task2.Handler, opsrepo.OperationsProvider, vo.ID, string) {
	t.Helper()
	h, p := newHandler()
	farmID := vo.NewID()
	ctx := testutil.OrgCtxWithID(farmID)

	res, err := h.Create(ctx, &task2.CreateTaskCommand{Title: "Задача"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	type idResult interface{ GetID() vo.ID }
	taskID := res.(response.IdResponse).ID.String()
	return h, p, farmID, taskID
}

// --- Start ---

func TestStartTask_ChangesStatusToInProgress(t *testing.T) {
	h, p, farmID, taskID := createAndGet(t)
	ctx := testutil.OrgCtxWithID(farmID)

	_, err := h.Start(ctx, &task2.TaskIDCommand{TaskID: taskID})
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	tasks, _ := p.Tasks().List(ctx, opsrepo.TaskFilter{FarmID: farmID})
	if tasks[0].Status != task.InProgress {
		t.Errorf("status: got %s, want IN_PROGRESS", tasks[0].Status)
	}
}

func TestStartTask_FailsWithoutOrgID(t *testing.T) {
	h, _, _, taskID := createAndGet(t)
	_, err := h.Start(testutil.OrgCtx(), &task2.TaskIDCommand{TaskID: taskID})
	// другой org — задача не найдена
	if err == nil {
		t.Fatal("expected error: task not found in different org")
	}
}

// --- Complete ---

func TestCompleteTask_ChangesStatusToDone(t *testing.T) {
	h, p, farmID, taskID := createAndGet(t)
	ctx := testutil.OrgCtxWithID(farmID)

	h.Complete(ctx, &task2.TaskIDCommand{TaskID: taskID})

	tasks, _ := p.Tasks().List(ctx, opsrepo.TaskFilter{FarmID: farmID})
	if tasks[0].Status != task.Done {
		t.Errorf("status: got %s, want DONE", tasks[0].Status)
	}
}

func TestCompleteTask_SetsCompletedAt(t *testing.T) {
	h, p, farmID, taskID := createAndGet(t)
	ctx := testutil.OrgCtxWithID(farmID)

	h.Complete(ctx, &task2.TaskIDCommand{TaskID: taskID})

	tasks, _ := p.Tasks().List(ctx, opsrepo.TaskFilter{FarmID: farmID})
	if tasks[0].CompletedAt == nil {
		t.Error("CompletedAt should be set after Complete")
	}
}

// --- Cancel ---

func TestCancelTask_ChangesStatusToCancelled(t *testing.T) {
	h, p, farmID, taskID := createAndGet(t)
	ctx := testutil.OrgCtxWithID(farmID)

	h.Cancel(ctx, &task2.TaskIDCommand{TaskID: taskID})

	tasks, _ := p.Tasks().List(ctx, opsrepo.TaskFilter{FarmID: farmID})
	if tasks[0].Status != task.Cancelled {
		t.Errorf("status: got %s, want CANCELLED", tasks[0].Status)
	}
}

// --- Assign ---

func TestAssignTask_SetsAssignment(t *testing.T) {
	h, p, farmID, taskID := createAndGet(t)
	ctx := testutil.OrgCtxWithID(farmID)
	userID := vo.NewID()

	_, err := h.Assign(ctx, &task2.AssignTaskCommand{
		TaskID: taskID,
		UserID: userID.String(),
	})
	if err != nil {
		t.Fatalf("assign: %v", err)
	}

	tasks, _ := p.Tasks().List(ctx, opsrepo.TaskFilter{FarmID: farmID})
	if tasks[0].Assignment == nil {
		t.Fatal("Assignment should be set")
	}
	if tasks[0].Assignment.UserID != userID {
		t.Errorf("userID: got %s, want %s", tasks[0].Assignment.UserID, userID)
	}
}

func TestAssignTask_FailsWithWrongPayloadType(t *testing.T) {
	h, _, _, _ := createAndGet(t)
	_, err := h.Assign(testutil.OrgCtx(), "bad payload")
	if err == nil {
		t.Fatal("expected error for wrong payload type")
	}
}
