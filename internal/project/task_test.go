package project

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	t.Run("valid task", func(t *testing.T) {
		id := TaskID("123")
		title := "Test Task"
		description := "Test Description"
		dueDate := time.Now().Add(24 * time.Hour)
		assigneeName := "John Doe"

		task, err := NewTask(id, title, &description, &dueDate, &assigneeName)
		require.NoError(t, err)
		assert.Equal(t, id, task.ID())
		assert.Equal(t, title, task.Title())
		assert.Equal(t, &description, task.Description())
		assert.Equal(t, &dueDate, task.DueDate())
		assert.Equal(t, &assigneeName, task.Asignee())
		assert.Equal(t, TaskStatusPending, task.Status())
	})

	t.Run("title too long", func(t *testing.T) {
		id := TaskID("123")
		title := "This is a very long title that exceeds the maximum length allowed"
		description := "Test Description"
		dueDate := time.Now().Add(24 * time.Hour)
		assigneeName := "John Doe"

		task, err := NewTask(id, title, &description, &dueDate, &assigneeName)
		assert.ErrorIs(t, err, ErrTitleTooLong)
		assert.Nil(t, task)
	})

	t.Run("due date in past", func(t *testing.T) {
		id := TaskID("123")
		title := "Test Task"
		description := "Test Description"
		dueDate := time.Now().Add(-24 * time.Hour)
		assigneeName := "John Doe"

		task, err := NewTask(id, title, &description, &dueDate, &assigneeName)
		assert.ErrorIs(t, err, ErrDueDateInPast)
		assert.Nil(t, task)
	})
}

func TestTaskChangeStatus(t *testing.T) {
	t.Run("valid status change", func(t *testing.T) {
		task := createValidTask(t)
		err := task.ChangeStatus(TaskStatusInProgress)
		require.NoError(t, err)
		assert.Equal(t, TaskStatusInProgress, task.Status())
		assert.Nil(t, task.CompletedAt())
	})

	t.Run("invalid status", func(t *testing.T) {
		task := createValidTask(t)
		err := task.ChangeStatus("invalid")
		assert.ErrorIs(t, err, ErrInvalidStatus)
	})

	t.Run("completed status sets completedAt", func(t *testing.T) {
		task := createValidTask(t)
		err := task.ChangeStatus(TaskStatusCompleted)
		require.NoError(t, err)
		assert.Equal(t, TaskStatusCompleted, task.Status())
		assert.NotNil(t, task.CompletedAt())
	})

	t.Run("changing from completed clears completedAt", func(t *testing.T) {
		task := createValidTask(t)
		_ = task.ChangeStatus(TaskStatusCompleted)
		err := task.ChangeStatus(TaskStatusInProgress)
		require.NoError(t, err)
		assert.Equal(t, TaskStatusInProgress, task.Status())
		assert.Nil(t, task.CompletedAt())
	})
}

func TestTaskEdit(t *testing.T) {
	t.Run("edit all fields", func(t *testing.T) {
		task := createValidTask(t)
		newTitle := "Updated Title"
		newDesc := "Updated Description"
		newDueDate := time.Now().Add(48 * time.Hour)
		newAssignee := "Jane Doe"
		newStatus := TaskStatusInProgress

		err := task.Edit(&newTitle, &newDesc, &newDueDate, &newAssignee, &newStatus)
		require.NoError(t, err)

		assert.Equal(t, newTitle, task.Title())
		assert.Equal(t, &newDesc, task.Description())
		assert.Equal(t, &newDueDate, task.DueDate())
		assert.Equal(t, &newAssignee, task.Asignee())
		assert.Equal(t, newStatus, task.Status())
	})

	t.Run("edit partial fields", func(t *testing.T) {
		task := createValidTask(t)
		originalDesc := task.Description()
		originalDueDate := task.DueDate()
		originalAssignee := task.Asignee()
		originalStatus := task.Status()

		newTitle := "Updated Title"
		err := task.Edit(&newTitle, nil, nil, nil, nil)
		require.NoError(t, err)

		assert.Equal(t, newTitle, task.Title())
		assert.Equal(t, originalDesc, task.Description())
		assert.Equal(t, originalDueDate, task.DueDate())
		assert.Equal(t, originalAssignee, task.Asignee())
		assert.Equal(t, originalStatus, task.Status())
	})

	t.Run("edit with invalid title", func(t *testing.T) {
		task := createValidTask(t)
		longTitle := "This is a very long title that exceeds the maximum length allowed"
		err := task.Edit(&longTitle, nil, nil, nil, nil)
		assert.ErrorIs(t, err, ErrTitleTooLong)
	})

	t.Run("edit with past due date", func(t *testing.T) {
		task := createValidTask(t)
		pastDate := time.Now().Add(-24 * time.Hour)
		err := task.Edit(nil, nil, &pastDate, nil, nil)
		assert.ErrorIs(t, err, ErrDueDateInPast)
	})

	t.Run("edit with invalid status", func(t *testing.T) {
		task := createValidTask(t)
		invalidStatus := TaskStatus("invalid")
		err := task.Edit(nil, nil, nil, nil, &invalidStatus)
		assert.ErrorIs(t, err, ErrInvalidStatus)
	})
}

func createValidTask(t *testing.T) *Task {
	id := TaskID("123")
	title := "Test Task"
	description := "Test Description"
	dueDate := time.Now().Add(24 * time.Hour)
	assigneeName := "John Doe"

	task, err := NewTask(id, title, &description, &dueDate, &assigneeName)
	require.NoError(t, err)
	return task
}
