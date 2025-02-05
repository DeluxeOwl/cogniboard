package adapters

import (
	"context"
	"fmt"
	"testing"
	"time"

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
		db:   "testdb",
		user: "testuser",
		pass: "testpass",
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

	db, err := ConnectDB(pgCreds.DSN(endpoint))
	require.NoError(t, err)
	require.NotNil(t, db)

	repo, err = NewPostgresTaskRepository(db)
	require.NoError(t, err)

	return repo, func() {
		testcontainers.CleanupContainer(t, postgresContainer)
	}
}

func createTaskWithID(t *testing.T, title string, description *string, dueDate *time.Time, assigneeID *project.TeamMemberID) (task *project.Task) {
	taskID, err := project.NewTaskID()
	require.NoError(t, err)

	task, err = project.NewTask(taskID, title, description, dueDate, assigneeID)
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
