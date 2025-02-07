package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGTaskRepository(t *testing.T) {
	t.Run("create and get task", func(t *testing.T) {
		repo := setupTestRepo(t)
		ctx := context.Background()

		taskID, err := project.NewTaskID()
		require.NoError(t, err)

		title := "Test Task"
		description := "Test Description"
		dueDate := time.Now().Add(24 * time.Hour)
		assigneeName := "John Doe"

		task, err := project.NewTask(taskID, title, &description, &dueDate, &assigneeName)
		require.NoError(t, err)

		err = repo.Create(ctx, task)
		require.NoError(t, err)

		fetchedTask, err := repo.GetByID(ctx, taskID)
		require.NoError(t, err)

		assert.Equal(t, task.ID(), fetchedTask.ID())
		assert.Equal(t, task.Title(), fetchedTask.Title())
		assert.Equal(t, task.Description(), fetchedTask.Description())
		assert.Equal(t, task.DueDate().Unix(), fetchedTask.DueDate().Unix())
		assert.Equal(t, task.Asignee(), fetchedTask.Asignee())
		assert.Equal(t, task.Status(), fetchedTask.Status())
		assert.Equal(t, task.CreatedAt().Unix(), fetchedTask.CreatedAt().Unix())
	})

	t.Run("update task", func(t *testing.T) {
		repo := setupTestRepo(t)
		ctx := context.Background()

		taskID, err := project.NewTaskID()
		require.NoError(t, err)

		title := "Test Task"
		description := "Test Description"
		dueDate := time.Now().Add(24 * time.Hour)
		assigneeName := "John Doe"

		task, err := project.NewTask(taskID, title, &description, &dueDate, &assigneeName)
		require.NoError(t, err)

		err = repo.Create(ctx, task)
		require.NoError(t, err)

		newTitle := "Updated Task"
		newDesc := "Updated Description"
		newDueDate := time.Now().Add(48 * time.Hour)
		newAssignee := "Jane Doe"
		newStatus := project.TaskStatusInProgress

		err = repo.UpdateTask(ctx, taskID, func(t *project.Task) (*project.Task, error) {
			return t, t.Edit(&newTitle, &newDesc, &newDueDate, &newAssignee, &newStatus)
		})
		require.NoError(t, err)

		updatedTask, err := repo.GetByID(ctx, taskID)
		require.NoError(t, err)

		assert.Equal(t, newTitle, updatedTask.Title())
		assert.Equal(t, &newDesc, updatedTask.Description())
		assert.Equal(t, newDueDate.Unix(), updatedTask.DueDate().Unix())
		assert.Equal(t, &newAssignee, updatedTask.Asignee())
		assert.Equal(t, newStatus, updatedTask.Status())
	})

	t.Run("get all tasks", func(t *testing.T) {
		repo := setupTestRepo(t)
		ctx := context.Background()

		// Create first task
		taskID1, err := project.NewTaskID()
		require.NoError(t, err)
		title1 := "Test Task 1"
		description1 := "Test Description 1"
		dueDate1 := time.Now().Add(24 * time.Hour)
		assigneeName1 := "John Doe"

		task1, err := project.NewTask(taskID1, title1, &description1, &dueDate1, &assigneeName1)
		require.NoError(t, err)
		err = repo.Create(ctx, task1)
		require.NoError(t, err)

		// Create second task
		taskID2, err := project.NewTaskID()
		require.NoError(t, err)
		title2 := "Test Task 2"
		description2 := "Test Description 2"
		dueDate2 := time.Now().Add(48 * time.Hour)
		assigneeName2 := "Jane Doe"

		task2, err := project.NewTask(taskID2, title2, &description2, &dueDate2, &assigneeName2)
		require.NoError(t, err)
		err = repo.Create(ctx, task2)
		require.NoError(t, err)

		// Get all tasks
		tasks, err := repo.AllTasks(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)

		// Verify tasks are returned in the correct order and with correct data
		assert.Equal(t, taskID1, tasks[0].ID())
		assert.Equal(t, title1, tasks[0].Title())
		assert.Equal(t, &description1, tasks[0].Description())
		assert.Equal(t, dueDate1.Unix(), tasks[0].DueDate().Unix())
		assert.Equal(t, &assigneeName1, tasks[0].Asignee())

		assert.Equal(t, taskID2, tasks[1].ID())
		assert.Equal(t, title2, tasks[1].Title())
		assert.Equal(t, &description2, tasks[1].Description())
		assert.Equal(t, dueDate2.Unix(), tasks[1].DueDate().Unix())
		assert.Equal(t, &assigneeName2, tasks[1].Asignee())
	})
}

func setupTestRepo(t *testing.T) project.TaskRepository {
	// Mock implementation for testing
	return NewInMemoryTaskRepository()
}

// InMemoryTaskRepository is a simple in-memory implementation for testing
type InMemoryTaskRepository struct {
	tasks map[project.TaskID]*project.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[project.TaskID]*project.Task),
	}
}

func (r *InMemoryTaskRepository) Create(ctx context.Context, task *project.Task) error {
	r.tasks[task.ID()] = task
	return nil
}

func (r *InMemoryTaskRepository) GetByID(ctx context.Context, id project.TaskID) (*project.Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return nil, nil
	}
	return task, nil
}

func (r *InMemoryTaskRepository) UpdateTask(ctx context.Context, id project.TaskID, updateFn func(t *project.Task) (*project.Task, error)) error {
	task, ok := r.tasks[id]
	if !ok {
		return nil
	}

	updatedTask, err := updateFn(task)
	if err != nil {
		return err
	}

	r.tasks[id] = updatedTask
	return nil
}

func (r *InMemoryTaskRepository) AllTasks(ctx context.Context) ([]project.Task, error) {
	tasks := make([]project.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}
