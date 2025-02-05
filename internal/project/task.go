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
	id          TaskID
	title       string
	description *string
	dueDate     *time.Time
	asigneeName *string
	createdAt   time.Time
	completedAt *time.Time
	status      TaskStatus
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

	task := &Task{
		id:          id,
		createdAt:   time.Now(),
		dueDate:     dueDate,
		asigneeName: assigneeName,
		title:       title,
		description: description,
		status:      TaskStatusPending,
	}
	return task, nil
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
	return nil
}

func (t *Task) AssignTo(memberName *string) error {
	if memberName == nil {
		return errors.New("user ID cannot be nil")
	}
	t.asigneeName = memberName
	return nil
}

func (t *Task) Unassign() {
	t.asigneeName = nil
}

func (t *Task) ID() TaskID {
	return t.id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() *string {
	if t.description == nil {
		return nil
	}
	copy := *t.description
	return &copy
}

func (t *Task) DueDate() *time.Time {
	if t.dueDate == nil {
		return nil
	}
	copy := *t.dueDate
	return &copy
}

func (t *Task) Asignee() *string {
	if t.asigneeName == nil {
		return nil
	}
	copy := *t.asigneeName
	return &copy
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) CompletedAt() *time.Time {
	if t.completedAt == nil {
		return nil
	}
	copy := *t.completedAt
	return &copy
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func UnmarshalTaskFromDB(
	id TaskID,
	title string,
	description *string,
	dueDate *time.Time,
	assigneeName *string,
	createdAt time.Time,
	completedAt *time.Time,
	status TaskStatus,
) (*Task, error) {
	task, err := NewTask(id, title, description, dueDate, assigneeName)
	if err != nil {
		return nil, err
	}

	task.createdAt = createdAt
	task.completedAt = completedAt
	task.status = status

	return task, nil
}
