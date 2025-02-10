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
		assert.Equal(t, id, task.id)
		assert.Equal(t, title, task.title)
		assert.Equal(t, &description, task.description)
		assert.Equal(t, &dueDate, task.dueDate)
		assert.Equal(t, &assigneeName, task.assigneeName)
		assert.Equal(t, TaskStatusPending, task.status)
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
		assert.Equal(t, TaskStatusInProgress, task.status)
		assert.Nil(t, task.completedAt)
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
		assert.Equal(t, TaskStatusCompleted, task.status)
		assert.NotNil(t, task.completedAt)
	})

	t.Run("changing from completed clears completedAt", func(t *testing.T) {
		task := createValidTask(t)
		_ = task.ChangeStatus(TaskStatusCompleted)
		err := task.ChangeStatus(TaskStatusInProgress)
		require.NoError(t, err)
		assert.Equal(t, TaskStatusInProgress, task.status)
		assert.Nil(t, task.completedAt)
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

		assert.Equal(t, newTitle, task.title)
		assert.Equal(t, &newDesc, task.description)
		assert.Equal(t, &newDueDate, task.dueDate)
		assert.Equal(t, &newAssignee, task.assigneeName)
		assert.Equal(t, newStatus, task.status)
	})

	t.Run("edit partial fields", func(t *testing.T) {
		task := createValidTask(t)
		originalDesc := task.description
		originalDueDate := task.dueDate
		originalAssignee := task.assigneeName
		originalStatus := task.status

		newTitle := "Updated Title"
		err := task.Edit(&newTitle, nil, nil, nil, nil)
		require.NoError(t, err)

		assert.Equal(t, newTitle, task.title)
		assert.Equal(t, originalDesc, task.description)
		assert.Equal(t, originalDueDate, task.dueDate)
		assert.Equal(t, originalAssignee, task.assigneeName)
		assert.Equal(t, originalStatus, task.status)
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
