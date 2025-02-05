package project

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	t.Run("creates task with valid input", func(t *testing.T) {
		// Arrange
		id := TaskID("test-id")
		title := "Test Task"
		description := "Test Description"
		now := time.Now().Add(24 * time.Hour)
		assigneeName := "john"

		// Act
		task, err := NewTask(id, title, &description, &now, &assigneeName)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, task.ID())
		assert.Equal(t, title, task.Title())
		assert.Equal(t, &description, task.Description())
		assert.Equal(t, &now, task.DueDate())
		assert.Equal(t, &assigneeName, task.Asignee())
		assert.Equal(t, TaskStatusPending, task.Status())
		assert.Nil(t, task.CompletedAt())
	})

	t.Run("fails with too long title", func(t *testing.T) {
		// Arrange
		id := TaskID("test-id")
		title := "This is a very long title that exceeds the maximum length limit"

		// Act
		task, err := NewTask(id, title, nil, nil, nil)

		// Assert
		assert.ErrorIs(t, err, ErrTitleTooLong)
		assert.Nil(t, task)
	})

	t.Run("fails with past due date", func(t *testing.T) {
		// Arrange
		id := TaskID("test-id")
		title := "Test Task"
		pastDate := time.Now().Add(-24 * time.Hour)

		// Act
		task, err := NewTask(id, title, nil, &pastDate, nil)

		// Assert
		assert.ErrorIs(t, err, ErrDueDateInPast)
		assert.Nil(t, task)
	})
}

func TestTaskStatusTransitions(t *testing.T) {
	// Arrange
	task, err := NewTask(TaskID("test-id"), "Test Task", nil, nil, nil)
	require.NoError(t, err)

	t.Run("completes task", func(t *testing.T) {
		// Act
		err := task.ChangeStatus(TaskStatusCompleted)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, TaskStatusCompleted, task.Status())
		assert.NotNil(t, task.CompletedAt())
	})

	t.Run("marks task as pending", func(t *testing.T) {
		// Act
		err := task.ChangeStatus(TaskStatusPending)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, TaskStatusPending, task.Status())
		assert.Nil(t, task.CompletedAt())
	})

	t.Run("marks task as in progress", func(t *testing.T) {
		// Act
		err := task.ChangeStatus(TaskStatusInProgress)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, TaskStatusInProgress, task.Status())
		assert.Nil(t, task.CompletedAt())
	})

	t.Run("marks task as in review", func(t *testing.T) {
		// Act
		err := task.ChangeStatus(TaskStatusInReview)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, TaskStatusInReview, task.Status())
		assert.Nil(t, task.CompletedAt())
	})
}

func TestTaskAssignment(t *testing.T) {
	// Arrange
	task, err := NewTask(TaskID("test-id"), "Test Task", nil, nil, nil)
	require.NoError(t, err)

	t.Run("assigns task to user", func(t *testing.T) {
		// Arrange
		memberName := "john"

		// Act
		err := task.AssignTo(&memberName)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &memberName, task.Asignee())
	})

	t.Run("fails to assign with nil user ID", func(t *testing.T) {
		// Act
		err := task.AssignTo(nil)

		// Assert
		assert.Error(t, err)
	})

	t.Run("unassigns task", func(t *testing.T) {
		// Act
		task.Unassign()

		// Assert
		assert.Nil(t, task.Asignee())
	})
}

func TestTaskGetters(t *testing.T) {
	// Arrange
	description := "Test Description"
	dueDate := time.Now().Add(24 * time.Hour)
	assigneeName := "john"
	task, err := NewTask(TaskID("test-id"), "Test Task", &description, &dueDate, &assigneeName)
	require.NoError(t, err)

	t.Run("returns copies of pointer values", func(t *testing.T) {
		// Act & Assert
		descCopy := task.Description()
		require.NotNil(t, descCopy)
		*descCopy = "Modified"
		assert.Equal(t, "Test Description", *task.Description())

		dueDateCopy := task.DueDate()
		require.NotNil(t, dueDateCopy)
		*dueDateCopy = dueDateCopy.Add(time.Hour)
		assert.Equal(t, dueDate, *task.DueDate())

		assigneeCopy := task.Asignee()
		require.NotNil(t, assigneeCopy)
		*assigneeCopy = "mary"
		assert.Equal(t, assigneeName, *task.Asignee())
	})

	t.Run("handles nil values", func(t *testing.T) {
		// Arrange
		task, err := NewTask(TaskID("test-id"), "Test Task", nil, nil, nil)
		require.NoError(t, err)

		// Act & Assert
		assert.Nil(t, task.Description())
		assert.Nil(t, task.DueDate())
		assert.Nil(t, task.Asignee())
		assert.Nil(t, task.CompletedAt())
	})
}

func TestUnmarshalFromDB(t *testing.T) {
	t.Run("unmarshals valid task", func(t *testing.T) {
		// Arrange
		id := TaskID("test-id")
		title := "Test Task"
		description := "Test Description"
		now := time.Now()
		dueDate := now.Add(24 * time.Hour)
		assigneeName := "john"
		completedAt := now.Add(12 * time.Hour)

		// Act
		task, err := UnmarshalTaskFromDB(
			id,
			title,
			&description,
			&dueDate,
			&assigneeName,
			now,
			&completedAt,
			TaskStatusCompleted,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, task.ID())
		assert.Equal(t, title, task.Title())
		assert.Equal(t, &description, task.Description())
		assert.Equal(t, &dueDate, task.DueDate())
		assert.Equal(t, &assigneeName, task.Asignee())
		assert.Equal(t, now, task.CreatedAt())
		assert.Equal(t, &completedAt, task.CompletedAt())
		assert.Equal(t, TaskStatusCompleted, task.Status())
	})

	t.Run("fails with invalid task data", func(t *testing.T) {
		// Arrange
		title := "This is a very long title that exceeds the maximum length limit"
		now := time.Now()

		// Act
		task, err := UnmarshalTaskFromDB(
			TaskID("test-id"),
			title,
			nil,
			nil,
			nil,
			now,
			nil,
			TaskStatusPending,
		)

		// Assert
		assert.ErrorIs(t, err, ErrTitleTooLong)
		assert.Nil(t, task)
	})
}
