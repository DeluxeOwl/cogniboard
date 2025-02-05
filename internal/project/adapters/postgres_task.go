package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// postgres driver
	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/jmoiron/sqlx"
)

type taskDTO struct {
	ID           string     `db:"id"`
	Title        string     `db:"title"`
	Description  *string    `db:"description"`
	DueDate      *time.Time `db:"due_date"`
	AssigneeName *string    `db:"assignee"`
	CreatedAt    time.Time  `db:"created_at"`
	CompletedAt  *time.Time `db:"completed_at"`
	Status       string     `db:"status"`
}

func (dto *taskDTO) toDomain() (*project.Task, error) {
	var assigneeName *string
	if dto.AssigneeName != nil {
		id := string(*dto.AssigneeName)
		assigneeName = &id
	}

	return project.UnmarshalTaskFromDB(
		project.TaskID(dto.ID),
		dto.Title,
		dto.Description,
		dto.DueDate,
		assigneeName,
		dto.CreatedAt,
		dto.CompletedAt,
		project.TaskStatus(dto.Status),
	)
}

func toDTO(t *project.Task) *taskDTO {
	var assigneeName *string
	if t.Asignee() != nil {
		nr := string(*t.Asignee())
		assigneeName = &nr
	}

	return &taskDTO{
		ID:           string(t.ID()),
		Title:        t.Title(),
		Description:  t.Description(),
		DueDate:      t.DueDate(),
		AssigneeName: assigneeName,
		CreatedAt:    t.CreatedAt(),
		CompletedAt:  t.CompletedAt(),
		Status:       string(t.Status()),
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

	return repo, nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *project.Task) error {
	dto := toDTO(task)

	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO tasks (id, title, description, due_date, assignee, created_at, completed_at, status)
		VALUES (:id, :title, :description, :due_date, :assignee, :created_at, :completed_at, :status)`,
		dto,
	)

	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id project.TaskID) (*project.Task, error) {
	var dto taskDTO
	err := r.db.GetContext(ctx, &dto,
		`SELECT id, title, description, due_date, assignee, created_at, completed_at, status
		FROM tasks WHERE id = $1`,
		string(id),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found: %w", err)
		}
		return nil, fmt.Errorf("get task: %w", err)
	}

	return dto.toDomain()
}

func (r *PostgresTaskRepository) UpdateTask(ctx context.Context, id project.TaskID, updateFn func(t *project.Task) (*project.Task, error)) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var dto taskDTO
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

	existingTask, err := dto.toDomain()
	if err != nil {
		return fmt.Errorf("convert to domain: %w", err)
	}

	updatedTask, err := updateFn(existingTask)
	if err != nil {
		return err
	}

	updatedDTO := toDTO(updatedTask)
	_, err = tx.NamedExecContext(ctx,
		`UPDATE tasks 
		SET title = :title, description = :description, due_date = :due_date, 
			assignee = :assignee, created_at = :created_at, 
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
