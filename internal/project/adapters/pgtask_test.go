package adapters

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/postgres"
	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type pgTestCreds struct {
	db   string
	user string
	pass string
}

func (creds *pgTestCreds) DSN(endpoint string) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		creds.user,
		creds.pass,
		endpoint,
		creds.db,
	)
}

var (
	pgImage = "postgres:17-alpine"
	pgPort  = "5432/tcp"
	pgCreds = pgTestCreds{
		db:   "cogniboard",
		user: "cogniboard",
		pass: "password",
	}
)

func setupPostgresRepo(ctx context.Context, t *testing.T) (repo *PostgresTaskRepository, cleanup func()) {
	req := testcontainers.ContainerRequest{
		Image:        pgImage,
		ExposedPorts: []string{pgPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(pgPort)),
		Env: map[string]string{
			"POSTGRES_DB":       pgCreds.db,
			"POSTGRES_USER":     pgCreds.user,
			"POSTGRES_PASSWORD": pgCreds.pass,
		},
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	endpoint, err := postgresContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	db, err := postgres.NewPostgresWithMigrate(ctx, pgCreds.DSN(endpoint))
	require.NoError(t, err)
	require.NotNil(t, db)

	repo, err = NewPostgresTaskRepository(db)
	require.NoError(t, err)

	return repo, func() {
		testcontainers.CleanupContainer(t, postgresContainer)
	}
}

func createTaskWithID(t *testing.T, title string, description *string, dueDate *time.Time, assigneeName *string) (task *project.Task) {
	taskID, err := project.NewTaskID()
	require.NoError(t, err)

	task, err = project.NewTask(taskID, title, description, dueDate, assigneeName)
	require.NoError(t, err)

	return task
}

func Test_RepoCreate(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupPostgresRepo(ctx, t)
	defer cleanup()

	task := createTaskWithID(t, "some new task", nil, nil, nil)

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	taskFromDB, err := repo.GetByID(ctx, task.GetSnapshot().ID)
	require.NoError(t, err)

	require.Equal(t, taskFromDB.GetSnapshot().ID, task.GetSnapshot().ID)
}

func Test_RepoGetByID(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupPostgresRepo(ctx, t)
	defer cleanup()

	t.Run("returns error when task does not exist", func(t *testing.T) {
		taskID, err := project.NewTaskID()
		require.NoError(t, err)

		task, err := repo.GetByID(ctx, taskID)
		require.Error(t, err)
		require.Nil(t, task)
		require.Contains(t, err.Error(), "task not found")
	})

	t.Run("successfully retrieves existing task with all fields", func(t *testing.T) {
		description := "test description"
		dueDate := time.Now().Add(24 * time.Hour)
		assigneeName := "john"
		task := createTaskWithID(t, "test task", &description, &dueDate, &assigneeName)

		err := repo.Create(ctx, task)
		require.NoError(t, err)

		taskFromDB, err := repo.GetByID(ctx, task.GetSnapshot().ID)
		require.NoError(t, err)
		require.NotNil(t, taskFromDB)

		require.Equal(t, task.GetSnapshot().ID, taskFromDB.GetSnapshot().ID)
		require.Equal(t, task.GetSnapshot().Title, taskFromDB.GetSnapshot().Title)
		require.Equal(t, task.GetSnapshot().Description, taskFromDB.GetSnapshot().Description)
		require.Equal(t, task.GetSnapshot().DueDate.UTC(), taskFromDB.GetSnapshot().DueDate.UTC())
		require.Equal(t, task.GetSnapshot().Asignee, taskFromDB.GetSnapshot().Asignee)
		require.Equal(t, task.GetSnapshot().Status, taskFromDB.GetSnapshot().Status)
		// Truncate timestamps to milliseconds for comparison
		require.Equal(t, task.GetSnapshot().CreatedAt.UTC().Truncate(time.Millisecond),
			taskFromDB.GetSnapshot().CreatedAt.UTC().Truncate(time.Millisecond))
		require.Equal(t, task.GetSnapshot().CompletedAt, taskFromDB.GetSnapshot().CompletedAt)
		require.Equal(t, task.GetSnapshot().UpdatedAt.UTC().Truncate(time.Millisecond),
			taskFromDB.GetSnapshot().UpdatedAt.UTC().Truncate(time.Millisecond))
	})
}

func Test_RepoUpdateTask(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupPostgresRepo(ctx, t)
	defer cleanup()

	t.Run("returns error when task does not exist", func(t *testing.T) {
		taskID, err := project.NewTaskID()
		require.NoError(t, err)

		err = repo.UpdateTask(ctx, taskID, func(t *project.Task) (*project.Task, error) {
			return t, nil
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "task not found")
	})

	t.Run("successfully updates task status", func(t *testing.T) {
		task := createTaskWithID(t, "test task", nil, nil, nil)
		err := repo.Create(ctx, task)
		require.NoError(t, err)

		err = repo.UpdateTask(ctx, task.GetSnapshot().ID, func(t *project.Task) (*project.Task, error) {
			err := t.ChangeStatus(project.TaskStatusCompleted)
			if err != nil {
				return nil, err
			}
			return t, nil
		})
		require.NoError(t, err)

		updatedTask, err := repo.GetByID(ctx, task.GetSnapshot().ID)
		snap := updatedTask.GetSnapshot()
		require.NoError(t, err)
		require.Equal(t, project.TaskStatusCompleted, snap.Status)
		require.NotNil(t, snap.CompletedAt)
	})

	t.Run("rolls back transaction on update function error", func(t *testing.T) {
		task := createTaskWithID(t, "test task", nil, nil, nil)
		err := repo.Create(ctx, task)
		require.NoError(t, err)

		expectedError := fmt.Errorf("update error")
		err = repo.UpdateTask(ctx, task.GetSnapshot().ID, func(t *project.Task) (*project.Task, error) {
			err := t.ChangeStatus(project.TaskStatusCompleted)
			if err != nil {
				return nil, err
			}
			return t, expectedError
		})
		require.ErrorIs(t, err, expectedError)

		unchangedTask, err := repo.GetByID(ctx, task.GetSnapshot().ID)
		snap := unchangedTask.GetSnapshot()
		require.NoError(t, err)
		require.Equal(t, project.TaskStatusPending, snap.Status)
		require.Nil(t, snap.CompletedAt)
	})

	t.Run("successfully updates task fields and sets updated_at", func(t *testing.T) {
		description := "initial description"
		dueDate := time.Now().Add(24 * time.Hour)
		assignee := "initial assignee"
		task := createTaskWithID(t, "initial title", &description, &dueDate, &assignee)

		err := repo.Create(ctx, task)
		require.NoError(t, err)

		// Store initial timestamps for comparison
		initialCreatedAt := task.GetSnapshot().CreatedAt
		initialUpdatedAt := task.GetSnapshot().UpdatedAt

		// Wait a moment to ensure updated_at will be different
		time.Sleep(time.Millisecond * 10)

		newTitle := "updated title"
		newDescription := "updated description"
		newDueDate := time.Now().Add(48 * time.Hour)
		newAssignee := "updated assignee"

		err = repo.UpdateTask(ctx, task.GetSnapshot().ID, func(t *project.Task) (*project.Task, error) {
			newStatus := project.TaskStatusInProgress
			err := t.Edit(&newTitle, &newDescription, &newDueDate, &newAssignee, &newStatus)
			if err != nil {
				return nil, err
			}
			return t, nil
		})
		require.NoError(t, err)

		updatedTask, err := repo.GetByID(ctx, task.GetSnapshot().ID)
		require.NoError(t, err)

		snap := updatedTask.GetSnapshot()
		// Verify all fields were updated
		require.Equal(t, newTitle, snap.Title)
		require.Equal(t, &newDescription, snap.Description)
		require.Equal(t, newDueDate.UTC(), snap.DueDate.UTC())
		require.Equal(t, &newAssignee, snap.Asignee)
		require.Equal(t, project.TaskStatusInProgress, snap.Status)

		// Verify timestamps
		require.Equal(t, initialCreatedAt.UTC(), snap.CreatedAt.UTC())
		require.True(t, snap.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("successfully updates individual fields", func(t *testing.T) {
		task := createTaskWithID(t, "test task", nil, nil, nil)
		err := repo.Create(ctx, task)
		require.NoError(t, err)

		// Update only title
		newTitle := "new title"
		err = repo.UpdateTask(ctx, task.GetSnapshot().ID, func(t *project.Task) (*project.Task, error) {
			err := t.Edit(&newTitle, nil, nil, nil, nil)
			if err != nil {
				return nil, err
			}
			return t, nil
		})
		require.NoError(t, err)

		updatedTask, err := repo.GetByID(ctx, task.GetSnapshot().ID)
		require.NoError(t, err)

		snap := updatedTask.GetSnapshot()
		require.Equal(t, newTitle, snap.Title)
		require.Nil(t, snap.Description)
		require.Nil(t, snap.DueDate)
		require.Nil(t, snap.Asignee)
		require.Equal(t, project.TaskStatusPending, snap.Status) // Status remains unchanged

		// Update only description
		newDescription := "new description"
		err = repo.UpdateTask(ctx, snap.ID, func(t *project.Task) (*project.Task, error) {
			err := t.Edit(nil, &newDescription, nil, nil, nil)
			if err != nil {
				return nil, err
			}
			return t, nil
		})
		require.NoError(t, err)

		updatedTask, err = repo.GetByID(ctx, snap.ID)
		snap = updatedTask.GetSnapshot()
		require.NoError(t, err)
		require.Equal(t, newTitle, snap.Title) // Title remains unchanged
		require.Equal(t, &newDescription, snap.Description)
		require.Nil(t, snap.DueDate)
		require.Nil(t, snap.Asignee)
	})
}
