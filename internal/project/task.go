package project

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	UpdateTask(ctx context.Context, id TaskID, updateFn func(t *Task) (*Task, error)) error
	AllTasks(ctx context.Context) ([]Task, error)
	AddFiles(ctx context.Context, taskID TaskID, files []File) error
}

type TaskID string

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusInReview   TaskStatus = "in_review"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	id           TaskID
	title        string
	description  *string
	dueDate      *time.Time
	assigneeName *string
	createdAt    time.Time
	updatedAt    time.Time
	completedAt  *time.Time
	status       TaskStatus
	files        []File
}

func NewTaskID() (TaskID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return TaskID(id.String()), nil
}

const MaxTitleLength = 50

var (
	ErrTitleTooLong  = fmt.Errorf("title cannot be longer than %d characters", MaxTitleLength)
	ErrDueDateInPast = errors.New("due date cannot be in the past")
	ErrInvalidStatus = errors.New("invalid task status")
)

var validTaskStatuses = map[TaskStatus]bool{
	TaskStatusPending:    true,
	TaskStatusInProgress: true,
	TaskStatusInReview:   true,
	TaskStatusCompleted:  true,
}

func NewTask(id TaskID, title string, description *string, dueDate *time.Time, assigneeName *string) (*Task, error) {
	if len(title) > 50 {
		return nil, ErrTitleTooLong
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		return nil, ErrDueDateInPast
	}

	now := time.Now()
	task := &Task{
		id:           id,
		createdAt:    now,
		updatedAt:    now,
		dueDate:      dueDate,
		assigneeName: assigneeName,
		title:        title,
		description:  description,
		status:       TaskStatusPending,
		files:        make([]File, 0),
	}
	return task, nil
}

func (t *Task) AddFile(file File) {
	t.files = append(t.files, file)
	t.updatedAt = time.Now()
}

func (t *Task) Files() []File {
	return t.files
}

func (t *Task) Edit(title *string, description *string, dueDate *time.Time, assigneeName *string, status *TaskStatus) error {
	if title != nil {
		if len(*title) > MaxTitleLength {
			return ErrTitleTooLong
		}
		t.title = *title
	}

	if dueDate != nil {
		if dueDate.Before(time.Now()) {
			return ErrDueDateInPast
		}
		t.dueDate = dueDate
	}

	if description != nil {
		t.description = description
	}

	if assigneeName != nil {
		t.assigneeName = assigneeName
	}

	if status != nil {
		if err := t.ChangeStatus(*status); err != nil {
			return err
		}
	}

	t.updatedAt = time.Now()
	return nil
}

func (t *Task) ChangeStatus(status TaskStatus) error {
	if !validTaskStatuses[status] {
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	if status == TaskStatusCompleted {
		t.completedAt = lo.ToPtr(time.Now())
	} else {
		t.completedAt = nil
	}
	t.status = status
	t.updatedAt = time.Now()
	return nil
}

type TaskSnapshot struct {
	ID          TaskID     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Assignee    *string    `json:"assignee"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      TaskStatus `json:"status"`
	Files       []File     `json:"files"`
}

// Used by the db adapters
func (t *Task) GetSnapshot() *TaskSnapshot {
	return &TaskSnapshot{
		ID:          t.id,
		Title:       t.title,
		Description: t.description,
		DueDate:     t.dueDate,
		Assignee:    t.assigneeName,
		CreatedAt:   t.createdAt,
		UpdatedAt:   t.updatedAt,
		CompletedAt: t.completedAt,
		Status:      t.status,
		Files:       t.files,
	}
}
