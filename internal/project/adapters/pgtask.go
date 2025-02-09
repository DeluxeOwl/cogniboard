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

func (r *PostgresTaskRepository) AddFiles(ctx context.Context, taskID project.TaskID, files []project.File) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First check if the task exists
	var exists bool
	err = tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`, string(taskID))
	if err != nil {
		return fmt.Errorf("check task existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("task not found: %w", sql.ErrNoRows)
	}

	// Insert files and create relationships
	for _, file := range files {
		// Insert file metadata
		fileID, err := project.NewTaskID() // Reuse TaskID type as it's also a UUID
		if err != nil {
			return fmt.Errorf("generate file ID: %w", err)
		}
		_, err = tx.ExecContext(ctx,
			`INSERT INTO files (id, name, size, mime_type, uploaded_at)
			VALUES ($1, $2, $3, $4, $5)`,
			string(fileID), file.Name, file.Size, file.MimeType, file.UploadedAt,
		)
		if err != nil {
			return fmt.Errorf("insert file: %w", err)
		}

		// Create task-file relationship
		_, err = tx.ExecContext(ctx,
			`INSERT INTO task_files (task_id, file_id)
			VALUES ($1, $2)`,
			string(taskID), string(fileID),
		)
		if err != nil {
			return fmt.Errorf("create task-file relationship: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
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
	var taskDTO project.OutTaskDTO
	err := r.db.GetContext(ctx, &taskDTO,
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

	// Get associated files
	var fileDTOs []project.OutFileDTO
	err = r.db.SelectContext(ctx, &fileDTOs,
		`SELECT f.id, f.name, f.size, f.mime_type, f.uploaded_at
		FROM files f
		JOIN task_files tf ON tf.file_id = f.id
		WHERE tf.task_id = $1`,
		string(id),
	)
	if err != nil {
		return nil, fmt.Errorf("get task files: %w", err)
	}

	task, err := project.FromOutTaskDTO(&taskDTO)
	if err != nil {
		return nil, fmt.Errorf("convert to domain: %w", err)
	}

	// Add files to task
	for _, fileDTO := range fileDTOs {
		task.AddFile(*project.FileFromOutFileDTO(&fileDTO))
	}

	return task, nil
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
	var taskDTOs []project.OutTaskDTO
	err := r.db.SelectContext(ctx, &taskDTOs,
		`SELECT id, title, description, due_date, assignee, created_at, updated_at, completed_at, status
		FROM tasks`,
	)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}

	tasks := make([]project.Task, 0, len(taskDTOs))
	for _, taskDTO := range taskDTOs {
		task, err := project.FromOutTaskDTO(&taskDTO)
		if err != nil {
			return nil, fmt.Errorf("convert task to domain: %w", err)
		}

		// Get files for each task
		var fileDTOs []project.OutFileDTO
		err = r.db.SelectContext(ctx, &fileDTOs,
			`SELECT f.id, f.name, f.size, f.mime_type, f.uploaded_at
			FROM files f
			JOIN task_files tf ON tf.file_id = f.id
			WHERE tf.task_id = $1`,
			taskDTO.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("get task files: %w", err)
		}

		for _, fileDTO := range fileDTOs {
			task.AddFile(*project.FileFromOutFileDTO(&fileDTO))
		}

		tasks = append(tasks, *task)
	}

	return tasks, nil
}
