package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// postgres driver
	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/jmoiron/sqlx"
)

type PostgresTaskRepository struct {
	db *sqlx.DB
}

var _ project.TaskRepository = &PostgresTaskRepository{}

func NewPostgresTaskRepository(db *sqlx.DB) (*PostgresTaskRepository, error) {
	if db == nil {
		return nil, errors.New("db connection cannot be nil")
	}

	repo := &PostgresTaskRepository{
		db: db,
	}

	return repo, nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *project.Task) error {
	dto := project.ToOutTaskDTO(task)

	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO tasks (id, title, description, due_date, assignee, created_at, updated_at, completed_at, status)
		VALUES (:id, :title, :description, :due_date, :assignee, :created_at, :updated_at, :completed_at, :status)`,
		dto,
	)

	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id project.TaskID) (*project.Task, error) {
	var dto project.OutTaskDTO
	err := r.db.GetContext(ctx, &dto,
		`SELECT id, title, description, due_date, assignee, created_at, updated_at, completed_at, status
		FROM tasks WHERE id = $1`,
		string(id),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found: %w", err)
		}
		return nil, fmt.Errorf("get task: %w", err)
	}

	return project.FromOutTaskDTO(&dto)
}

func (r *PostgresTaskRepository) UpdateTask(ctx context.Context, id project.TaskID, updateFn func(t *project.Task) (*project.Task, error)) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var dto project.OutTaskDTO
	err = tx.GetContext(ctx, &dto,
		`SELECT id, title, description, due_date, assignee, created_at, completed_at, status
		FROM tasks WHERE id = $1 FOR UPDATE`,
		string(id),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task not found: %w", err)
		}
		return fmt.Errorf("get task: %w", err)
	}

	existingTask, err := project.FromOutTaskDTO(&dto)
	if err != nil {
		return fmt.Errorf("convert to domain: %w", err)
	}

	updatedTask, err := updateFn(existingTask)
	if err != nil {
		return err
	}

	updatedDTO := project.ToOutTaskDTO(updatedTask)
	_, err = tx.NamedExecContext(ctx,
		`UPDATE tasks
		SET title = :title, description = :description, due_date = :due_date,
			assignee = :assignee, created_at = :created_at, updated_at = :updated_at,
			completed_at = :completed_at, status = :status
		WHERE id = :id`,
		updatedDTO,
	)

	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresTaskRepository) AllTasks(ctx context.Context) ([]project.Task, error) {
	var dtos []project.OutTaskDTO
	err := r.db.SelectContext(ctx, &dtos,
		`SELECT id, title, description, due_date, assignee, created_at, updated_at, completed_at, status
		FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}

	tasks := []project.Task{}
	for _, dto := range dtos {
		task, err := project.FromOutTaskDTO(&dto)
		if err != nil {
			return nil, fmt.Errorf("convert task to domain: %w", err)
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}
