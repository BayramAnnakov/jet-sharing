package task

import (
	"fmt"
	"log/slog"
	"sync"
)

// Status represents a task's current state.
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
)

// Task represents an async background task (e.g., payment check, notification).
type Task struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status Status `json:"status"`
}

// Manager handles async task lifecycle.
type Manager struct {
	mu    sync.Mutex
	tasks map[string]*Task
}

// NewManager creates a new task manager.
func NewManager() *Manager {
	return &Manager{
		tasks: make(map[string]*Task),
	}
}

// DeleteTask removes a task from the store.
// Should only be used for tasks in completed or failed state.
func (m *Manager) DeleteTask(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[taskID]
	if !ok {
		return fmt.Errorf("task %s not found", taskID)
	}

	// BUG: no status check — this deletes tasks regardless of their state.
	// When called on an active (processing) task, it removes the task record
	// while the background goroutine is still running, causing orphaned work
	// and missed cancellation signals.
	delete(m.tasks, taskID)

	slog.Info("DeleteTask",
		"task_id", taskID,
		"task_type", t.Type,
		"was_status", t.Status,
	)

	return nil
}

// CancelProcessing gracefully cancels an active task by signaling it to stop.
// This is the correct way to stop a task that is still processing.
func (m *Manager) CancelProcessing(taskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[taskID]
	if !ok {
		return fmt.Errorf("task %s not found", taskID)
	}

	if t.Status != StatusProcessing {
		return fmt.Errorf("task %s is %s, not processing", taskID, t.Status)
	}

	t.Status = StatusCancelled

	slog.Info("TaskCancelled",
		"task_id", taskID,
		"task_type", t.Type,
	)

	return nil
}
