package project

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type taskID string

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusInReview   TaskStatus = "in_review"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	id          taskID
	title       string
	description *string
	dueDate     *time.Time
	assigneeID  *TeamMemberID
	createdAt   time.Time
	completedAt *time.Time
	status      TaskStatus
}

func NewTaskID() (taskID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return taskID(id.String()), nil
}

const MaxTitleLength = 50

var (
	ErrTitleTooLong  = fmt.Errorf("title cannot be longer than %d characters", MaxTitleLength)
	ErrDueDateInPast = errors.New("due date cannot be in the past")
)

func NewTask(id taskID, title string, description *string, dueDate *time.Time, assigneeID *TeamMemberID) (*Task, error) {
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
		assigneeID:  assigneeID,
		title:       title,
		description: description,
		status:      TaskStatusPending,
	}
	return task, nil
}

func (t *Task) Complete() {
	t.completedAt = lo.ToPtr(time.Now())
	t.status = TaskStatusCompleted
}

func (t *Task) MarkPending() {
	t.completedAt = nil
	t.status = TaskStatusPending
}

func (t *Task) MarkInProgress() {
	t.completedAt = nil
	t.status = TaskStatusInProgress
}

func (t *Task) MarkInReview() {
	t.completedAt = nil
	t.status = TaskStatusInReview
}

func (t *Task) AssignTo(userID *TeamMemberID) error {
	if userID == nil {
		return errors.New("user ID cannot be nil")
	}
	t.assigneeID = userID
	return nil
}

func (t *Task) Unassign() {
	t.assigneeID = nil
}

func (t *Task) ID() taskID {
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

func (t *Task) AssigneeID() *TeamMemberID {
	if t.assigneeID == nil {
		return nil
	}
	copy := *t.assigneeID
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

func UnmarshalFromDB(
	id taskID,
	title string,
	description *string,
	dueDate *time.Time,
	assigneeID *TeamMemberID,
	createdAt time.Time,
	completedAt *time.Time,
	status TaskStatus,
) (*Task, error) {
	task, err := NewTask(id, title, description, dueDate, assigneeID)
	if err != nil {
		return nil, err
	}

	task.createdAt = createdAt
	task.completedAt = completedAt
	task.status = status

	return task, nil
}
