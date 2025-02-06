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

	db, err := postgres.NewPostgresWithMigrate(pgCreds.DSN(endpoint))
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

	taskFromDB, err := repo.GetByID(ctx, task.ID())
	require.NoError(t, err)

	require.Equal(t, taskFromDB.ID(), task.ID())
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

		taskFromDB, err := repo.GetByID(ctx, task.ID())
		require.NoError(t, err)
		require.NotNil(t, taskFromDB)

		require.Equal(t, task.ID(), taskFromDB.ID())
		require.Equal(t, task.Title(), taskFromDB.Title())
		require.Equal(t, task.Description(), taskFromDB.Description())
		require.Equal(t, task.DueDate().UTC(), taskFromDB.DueDate().UTC())
		require.Equal(t, task.Asignee(), taskFromDB.Asignee())
		require.Equal(t, task.Status(), taskFromDB.Status())
		require.Equal(t, task.CreatedAt().UTC(), taskFromDB.CreatedAt().UTC())
		require.Equal(t, task.CompletedAt(), taskFromDB.CompletedAt())
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

		err = repo.UpdateTask(ctx, task.ID(), func(t *project.Task) (*project.Task, error) {
			err := t.ChangeStatus(project.TaskStatusCompleted)
			if err != nil {
				return nil, err
			}
			return t, nil
		})
		require.NoError(t, err)

		updatedTask, err := repo.GetByID(ctx, task.ID())
		require.NoError(t, err)
		require.Equal(t, project.TaskStatusCompleted, updatedTask.Status())
		require.NotNil(t, updatedTask.CompletedAt())
	})

	t.Run("successfully updates task assignment", func(t *testing.T) {
		task := createTaskWithID(t, "test task", nil, nil, nil)
		err := repo.Create(ctx, task)
		require.NoError(t, err)

		assigneeName := "john"
		err = repo.UpdateTask(ctx, task.ID(), func(t *project.Task) (*project.Task, error) {
			return t, t.AssignTo(&assigneeName)
		})
		require.NoError(t, err)

		updatedTask, err := repo.GetByID(ctx, task.ID())
		require.NoError(t, err)
		require.Equal(t, &assigneeName, updatedTask.Asignee())
	})

	t.Run("rolls back transaction on update function error", func(t *testing.T) {
		task := createTaskWithID(t, "test task", nil, nil, nil)
		err := repo.Create(ctx, task)
		require.NoError(t, err)

		expectedError := fmt.Errorf("update error")
		err = repo.UpdateTask(ctx, task.ID(), func(t *project.Task) (*project.Task, error) {
			err := t.ChangeStatus(project.TaskStatusCompleted)
			if err != nil {
				return nil, err
			}
			return t, expectedError
		})
		require.ErrorIs(t, err, expectedError)

		unchangedTask, err := repo.GetByID(ctx, task.ID())
		require.NoError(t, err)
		require.Equal(t, project.TaskStatusPending, unchangedTask.Status())
		require.Nil(t, unchangedTask.CompletedAt())
	})
}
