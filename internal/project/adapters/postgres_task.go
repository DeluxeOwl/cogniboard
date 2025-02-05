package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	// postgres driver
	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/jmoiron/sqlx"
)

const createTasksTable = `
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    description TEXT,
    due_date TIMESTAMP WITH TIME ZONE,
    assignee_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL
);
`

type taskDTO struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	DueDate     *time.Time `db:"due_date"`
	AssigneeID  *string    `db:"assignee_id"`
	CreatedAt   time.Time  `db:"created_at"`
	CompletedAt *time.Time `db:"completed_at"`
	Status      string     `db:"status"`
}

func (dto *taskDTO) toDomain() (*project.Task, error) {
	var assigneeID *project.TeamMemberID
	if dto.AssigneeID != nil {
		id := project.TeamMemberID(*dto.AssigneeID)
		assigneeID = &id
	}

	return project.UnmarshalTaskFromDB(
		project.TaskID(dto.ID),
		dto.Title,
		dto.Description,
		dto.DueDate,
		assigneeID,
		dto.CreatedAt,
		dto.CompletedAt,
		project.TaskStatus(dto.Status),
	)
}

func toDTO(t *project.Task) *taskDTO {
	var assigneeID *string
	if t.AssigneeID() != nil {
		str := string(*t.AssigneeID())
		assigneeID = &str
	}

	return &taskDTO{
		ID:          string(t.ID()),
		Title:       t.Title(),
		Description: t.Description(),
		DueDate:     t.DueDate(),
		AssigneeID:  assigneeID,
		CreatedAt:   t.CreatedAt(),
		CompletedAt: t.CompletedAt(),
		Status:      string(t.Status()),
	}
}

type PostgresTaskRepository struct {
	db *sqlx.DB
}

func NewPostgresTaskRepository(db *sqlx.DB) (*PostgresTaskRepository, error) {
	if db == nil {
		return nil, errors.New("db connection cannot be nil")
	}

	repo := &PostgresTaskRepository{
		db: db,
	}

	if err := repo.createSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return repo, nil
}

func (r *PostgresTaskRepository) createSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, createTasksTable)
	return err
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *project.Task) error {
	dto := toDTO(task)

	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO tasks (id, title, description, due_date, assignee_id, created_at, completed_at, status)
		VALUES (:id, :title, :description, :due_date, :assignee_id, :created_at, :completed_at, :status)`,
		dto,
	)

	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id project.TaskID) (*project.Task, error) {
	var dto taskDTO
	err := r.db.GetContext(ctx, &dto,
		`SELECT id, title, description, due_date, assignee_id, created_at, completed_at, status
		FROM tasks WHERE id = $1`,
		string(id),
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return dto.toDomain()
}

func (r *PostgresTaskRepository) UpdateTask(ctx context.Context, id project.TaskID, updateFn func(t *project.Task) (*project.Task, error)) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var dto taskDTO
	err = tx.GetContext(ctx, &dto,
		`SELECT id, title, description, due_date, assignee_id, created_at, completed_at, status
		FROM tasks WHERE id = $1 FOR UPDATE`,
		string(id),
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("task not found: %w", err)
		}
		return fmt.Errorf("failed to get task: %w", err)
	}

	existingTask, err := dto.toDomain()
	if err != nil {
		return fmt.Errorf("failed to convert to domain: %w", err)
	}

	updatedTask, err := updateFn(existingTask)
	if err != nil {
		return err
	}

	updatedDTO := toDTO(updatedTask)
	_, err = tx.NamedExecContext(ctx,
		`UPDATE tasks 
		SET title = :title, description = :description, due_date = :due_date, 
			assignee_id = :assignee_id, created_at = :created_at, 
			completed_at = :completed_at, status = :status
		WHERE id = :id`,
		updatedDTO,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ConnectDB creates a new sqlx.DB instance with connection settings
func ConnectDB(dsn string) (*sqlx.DB, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Set timeouts
	config.ConnectTimeout = 10 * time.Second
	config.RuntimeParams["statement_timeout"] = "30000" // 30 seconds
	config.RuntimeParams["idle_in_transaction_session_timeout"] = "30000"

	connStr := stdlib.RegisterConnConfig(config)
	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
