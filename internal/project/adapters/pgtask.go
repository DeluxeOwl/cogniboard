package adapters

import (
	"context"
	"errors"
	"fmt"

	// postgres driver
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent"
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent/task"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type PostgresTaskRepository struct {
	client *ent.Client
}

var _ project.TaskRepository = &PostgresTaskRepository{}

func NewPostgresTaskRepository(client *ent.Client) (*PostgresTaskRepository, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}

	repo := &PostgresTaskRepository{
		client: client,
	}

	return repo, nil
}

func (r *PostgresTaskRepository) AddFiles(ctx context.Context, taskID project.TaskID, files []project.File) error {
	return WithTx(ctx, r.client, func(tx *ent.Tx) error {
		// First create all file records

		entFiles := make([]*ent.File, len(files))

		for i, f := range files {
			snap := f.GetSnapshot()
			file, err := tx.File.Create().
				SetID(snap.ID).
				SetName(snap.Name).
				SetSize(snap.Size).
				SetMimeType(snap.MimeType).
				SetUploadedAt(snap.UploadedAt).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
			entFiles[i] = file
		}

		// Then link the files to the task
		_, err := tx.Task.UpdateOneID(string(taskID)).
			AddFiles(entFiles...).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}

		return nil
	})
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *project.Task) error {
	t := task.GetSnapshot()

	_, err := r.client.Task.Create().
		SetID(string(t.ID)).
		SetTitle(t.Title).
		SetCreatedAt(t.CreatedAt).
		SetNillableAssigneeName(t.Assignee).
		SetNillableCompletedAt(t.CompletedAt).
		SetNillableDescription(t.Description).
		SetNillableDueDate(t.DueDate).
		Save(ctx)

	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id project.TaskID) (*project.Task, error) {
	task, err := r.client.Task.Query().
		Where(task.IDEQ(string(id))).
		WithFiles().
		First(ctx)
	if err != nil {
		return nil, err
	}

	return project.UnmarshalTaskFromDB(task)
}

func (r *PostgresTaskRepository) UpdateTask(ctx context.Context, id project.TaskID, updateFn func(t *project.Task) (*project.Task, error)) error {
	return WithTx(ctx, r.client, func(tx *ent.Tx) error {
		// Get existing task with files
		existingTask, err := tx.Task.Query().
			Where(task.IDEQ(string(id))).
			WithFiles().
			First(ctx)
		if err != nil {
			return fmt.Errorf("query task: %w", err)
		}

		// Convert to domain model
		domainTask, err := project.UnmarshalTaskFromDB(existingTask)
		if err != nil {
			return fmt.Errorf("convert to domain model: %w", err)
		}

		// Apply update function
		updatedTask, err := updateFn(domainTask)
		if err != nil {
			return fmt.Errorf("update function: %w", err)
		}

		snap := updatedTask.GetSnapshot()

		// Update in database
		_, err = tx.Task.UpdateOneID(string(id)).
			SetTitle(snap.Title).
			SetNillableDescription(snap.Description).
			SetNillableDueDate(snap.DueDate).
			SetNillableAssigneeName(snap.Assignee).
			SetNillableCompletedAt(snap.CompletedAt).
			SetStatus(task.Status(snap.Status)).
			SetUpdatedAt(snap.UpdatedAt).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("save task: %w", err)
		}

		return nil
	})
}

func (r *PostgresTaskRepository) AllTasks(ctx context.Context) ([]project.Task, error) {
	// Query all tasks with their files
	entTasks, err := r.client.Task.Query().
		WithFiles().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}

	// Convert to domain tasks
	tasks := make([]project.Task, 0, len(entTasks))
	for _, entTask := range entTasks {
		domainTask, err := project.UnmarshalTaskFromDB(entTask)
		if err != nil {
			return nil, fmt.Errorf("convert task %s: %w", entTask.ID, err)
		}
		tasks = append(tasks, *domainTask)
	}

	return tasks, nil
}

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
