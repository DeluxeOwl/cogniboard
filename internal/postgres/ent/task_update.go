// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent/file"
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent/predicate"
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent/task"
)

// TaskUpdate is the builder for updating Task entities.
type TaskUpdate struct {
	config
	hooks    []Hook
	mutation *TaskMutation
}

// Where appends a list predicates to the TaskUpdate builder.
func (tu *TaskUpdate) Where(ps ...predicate.Task) *TaskUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetTitle sets the "title" field.
func (tu *TaskUpdate) SetTitle(s string) *TaskUpdate {
	tu.mutation.SetTitle(s)
	return tu
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableTitle(s *string) *TaskUpdate {
	if s != nil {
		tu.SetTitle(*s)
	}
	return tu
}

// SetDescription sets the "description" field.
func (tu *TaskUpdate) SetDescription(s string) *TaskUpdate {
	tu.mutation.SetDescription(s)
	return tu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableDescription(s *string) *TaskUpdate {
	if s != nil {
		tu.SetDescription(*s)
	}
	return tu
}

// ClearDescription clears the value of the "description" field.
func (tu *TaskUpdate) ClearDescription() *TaskUpdate {
	tu.mutation.ClearDescription()
	return tu
}

// SetDueDate sets the "due_date" field.
func (tu *TaskUpdate) SetDueDate(t time.Time) *TaskUpdate {
	tu.mutation.SetDueDate(t)
	return tu
}

// SetNillableDueDate sets the "due_date" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableDueDate(t *time.Time) *TaskUpdate {
	if t != nil {
		tu.SetDueDate(*t)
	}
	return tu
}

// ClearDueDate clears the value of the "due_date" field.
func (tu *TaskUpdate) ClearDueDate() *TaskUpdate {
	tu.mutation.ClearDueDate()
	return tu
}

// SetAssigneeName sets the "assignee_name" field.
func (tu *TaskUpdate) SetAssigneeName(s string) *TaskUpdate {
	tu.mutation.SetAssigneeName(s)
	return tu
}

// SetNillableAssigneeName sets the "assignee_name" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableAssigneeName(s *string) *TaskUpdate {
	if s != nil {
		tu.SetAssigneeName(*s)
	}
	return tu
}

// ClearAssigneeName clears the value of the "assignee_name" field.
func (tu *TaskUpdate) ClearAssigneeName() *TaskUpdate {
	tu.mutation.ClearAssigneeName()
	return tu
}

// SetCreatedAt sets the "created_at" field.
func (tu *TaskUpdate) SetCreatedAt(t time.Time) *TaskUpdate {
	tu.mutation.SetCreatedAt(t)
	return tu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableCreatedAt(t *time.Time) *TaskUpdate {
	if t != nil {
		tu.SetCreatedAt(*t)
	}
	return tu
}

// SetUpdatedAt sets the "updated_at" field.
func (tu *TaskUpdate) SetUpdatedAt(t time.Time) *TaskUpdate {
	tu.mutation.SetUpdatedAt(t)
	return tu
}

// SetCompletedAt sets the "completed_at" field.
func (tu *TaskUpdate) SetCompletedAt(t time.Time) *TaskUpdate {
	tu.mutation.SetCompletedAt(t)
	return tu
}

// SetNillableCompletedAt sets the "completed_at" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableCompletedAt(t *time.Time) *TaskUpdate {
	if t != nil {
		tu.SetCompletedAt(*t)
	}
	return tu
}

// ClearCompletedAt clears the value of the "completed_at" field.
func (tu *TaskUpdate) ClearCompletedAt() *TaskUpdate {
	tu.mutation.ClearCompletedAt()
	return tu
}

// SetStatus sets the "status" field.
func (tu *TaskUpdate) SetStatus(t task.Status) *TaskUpdate {
	tu.mutation.SetStatus(t)
	return tu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableStatus(t *task.Status) *TaskUpdate {
	if t != nil {
		tu.SetStatus(*t)
	}
	return tu
}

// AddFileIDs adds the "files" edge to the File entity by IDs.
func (tu *TaskUpdate) AddFileIDs(ids ...string) *TaskUpdate {
	tu.mutation.AddFileIDs(ids...)
	return tu
}

// AddFiles adds the "files" edges to the File entity.
func (tu *TaskUpdate) AddFiles(f ...*File) *TaskUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return tu.AddFileIDs(ids...)
}

// Mutation returns the TaskMutation object of the builder.
func (tu *TaskUpdate) Mutation() *TaskMutation {
	return tu.mutation
}

// ClearFiles clears all "files" edges to the File entity.
func (tu *TaskUpdate) ClearFiles() *TaskUpdate {
	tu.mutation.ClearFiles()
	return tu
}

// RemoveFileIDs removes the "files" edge to File entities by IDs.
func (tu *TaskUpdate) RemoveFileIDs(ids ...string) *TaskUpdate {
	tu.mutation.RemoveFileIDs(ids...)
	return tu
}

// RemoveFiles removes "files" edges to File entities.
func (tu *TaskUpdate) RemoveFiles(f ...*File) *TaskUpdate {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return tu.RemoveFileIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TaskUpdate) Save(ctx context.Context) (int, error) {
	tu.defaults()
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TaskUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TaskUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TaskUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tu *TaskUpdate) defaults() {
	if _, ok := tu.mutation.UpdatedAt(); !ok {
		v := task.UpdateDefaultUpdatedAt()
		tu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tu *TaskUpdate) check() error {
	if v, ok := tu.mutation.Status(); ok {
		if err := task.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Task.status": %w`, err)}
		}
	}
	return nil
}

func (tu *TaskUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(task.Table, task.Columns, sqlgraph.NewFieldSpec(task.FieldID, field.TypeString))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.Title(); ok {
		_spec.SetField(task.FieldTitle, field.TypeString, value)
	}
	if value, ok := tu.mutation.Description(); ok {
		_spec.SetField(task.FieldDescription, field.TypeString, value)
	}
	if tu.mutation.DescriptionCleared() {
		_spec.ClearField(task.FieldDescription, field.TypeString)
	}
	if value, ok := tu.mutation.DueDate(); ok {
		_spec.SetField(task.FieldDueDate, field.TypeTime, value)
	}
	if tu.mutation.DueDateCleared() {
		_spec.ClearField(task.FieldDueDate, field.TypeTime)
	}
	if value, ok := tu.mutation.AssigneeName(); ok {
		_spec.SetField(task.FieldAssigneeName, field.TypeString, value)
	}
	if tu.mutation.AssigneeNameCleared() {
		_spec.ClearField(task.FieldAssigneeName, field.TypeString)
	}
	if value, ok := tu.mutation.CreatedAt(); ok {
		_spec.SetField(task.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.UpdatedAt(); ok {
		_spec.SetField(task.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.CompletedAt(); ok {
		_spec.SetField(task.FieldCompletedAt, field.TypeTime, value)
	}
	if tu.mutation.CompletedAtCleared() {
		_spec.ClearField(task.FieldCompletedAt, field.TypeTime)
	}
	if value, ok := tu.mutation.Status(); ok {
		_spec.SetField(task.FieldStatus, field.TypeEnum, value)
	}
	if tu.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   task.FilesTable,
			Columns: task.FilesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedFilesIDs(); len(nodes) > 0 && !tu.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   task.FilesTable,
			Columns: task.FilesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   task.FilesTable,
			Columns: task.FilesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{task.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TaskUpdateOne is the builder for updating a single Task entity.
type TaskUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TaskMutation
}

// SetTitle sets the "title" field.
func (tuo *TaskUpdateOne) SetTitle(s string) *TaskUpdateOne {
	tuo.mutation.SetTitle(s)
	return tuo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableTitle(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetTitle(*s)
	}
	return tuo
}

// SetDescription sets the "description" field.
func (tuo *TaskUpdateOne) SetDescription(s string) *TaskUpdateOne {
	tuo.mutation.SetDescription(s)
	return tuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableDescription(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetDescription(*s)
	}
	return tuo
}

// ClearDescription clears the value of the "description" field.
func (tuo *TaskUpdateOne) ClearDescription() *TaskUpdateOne {
	tuo.mutation.ClearDescription()
	return tuo
}

// SetDueDate sets the "due_date" field.
func (tuo *TaskUpdateOne) SetDueDate(t time.Time) *TaskUpdateOne {
	tuo.mutation.SetDueDate(t)
	return tuo
}

// SetNillableDueDate sets the "due_date" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableDueDate(t *time.Time) *TaskUpdateOne {
	if t != nil {
		tuo.SetDueDate(*t)
	}
	return tuo
}

// ClearDueDate clears the value of the "due_date" field.
func (tuo *TaskUpdateOne) ClearDueDate() *TaskUpdateOne {
	tuo.mutation.ClearDueDate()
	return tuo
}

// SetAssigneeName sets the "assignee_name" field.
func (tuo *TaskUpdateOne) SetAssigneeName(s string) *TaskUpdateOne {
	tuo.mutation.SetAssigneeName(s)
	return tuo
}

// SetNillableAssigneeName sets the "assignee_name" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableAssigneeName(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetAssigneeName(*s)
	}
	return tuo
}

// ClearAssigneeName clears the value of the "assignee_name" field.
func (tuo *TaskUpdateOne) ClearAssigneeName() *TaskUpdateOne {
	tuo.mutation.ClearAssigneeName()
	return tuo
}

// SetCreatedAt sets the "created_at" field.
func (tuo *TaskUpdateOne) SetCreatedAt(t time.Time) *TaskUpdateOne {
	tuo.mutation.SetCreatedAt(t)
	return tuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableCreatedAt(t *time.Time) *TaskUpdateOne {
	if t != nil {
		tuo.SetCreatedAt(*t)
	}
	return tuo
}

// SetUpdatedAt sets the "updated_at" field.
func (tuo *TaskUpdateOne) SetUpdatedAt(t time.Time) *TaskUpdateOne {
	tuo.mutation.SetUpdatedAt(t)
	return tuo
}

// SetCompletedAt sets the "completed_at" field.
func (tuo *TaskUpdateOne) SetCompletedAt(t time.Time) *TaskUpdateOne {
	tuo.mutation.SetCompletedAt(t)
	return tuo
}

// SetNillableCompletedAt sets the "completed_at" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableCompletedAt(t *time.Time) *TaskUpdateOne {
	if t != nil {
		tuo.SetCompletedAt(*t)
	}
	return tuo
}

// ClearCompletedAt clears the value of the "completed_at" field.
func (tuo *TaskUpdateOne) ClearCompletedAt() *TaskUpdateOne {
	tuo.mutation.ClearCompletedAt()
	return tuo
}

// SetStatus sets the "status" field.
func (tuo *TaskUpdateOne) SetStatus(t task.Status) *TaskUpdateOne {
	tuo.mutation.SetStatus(t)
	return tuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableStatus(t *task.Status) *TaskUpdateOne {
	if t != nil {
		tuo.SetStatus(*t)
	}
	return tuo
}

// AddFileIDs adds the "files" edge to the File entity by IDs.
func (tuo *TaskUpdateOne) AddFileIDs(ids ...string) *TaskUpdateOne {
	tuo.mutation.AddFileIDs(ids...)
	return tuo
}

// AddFiles adds the "files" edges to the File entity.
func (tuo *TaskUpdateOne) AddFiles(f ...*File) *TaskUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return tuo.AddFileIDs(ids...)
}

// Mutation returns the TaskMutation object of the builder.
func (tuo *TaskUpdateOne) Mutation() *TaskMutation {
	return tuo.mutation
}

// ClearFiles clears all "files" edges to the File entity.
func (tuo *TaskUpdateOne) ClearFiles() *TaskUpdateOne {
	tuo.mutation.ClearFiles()
	return tuo
}

// RemoveFileIDs removes the "files" edge to File entities by IDs.
func (tuo *TaskUpdateOne) RemoveFileIDs(ids ...string) *TaskUpdateOne {
	tuo.mutation.RemoveFileIDs(ids...)
	return tuo
}

// RemoveFiles removes "files" edges to File entities.
func (tuo *TaskUpdateOne) RemoveFiles(f ...*File) *TaskUpdateOne {
	ids := make([]string, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return tuo.RemoveFileIDs(ids...)
}

// Where appends a list predicates to the TaskUpdate builder.
func (tuo *TaskUpdateOne) Where(ps ...predicate.Task) *TaskUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TaskUpdateOne) Select(field string, fields ...string) *TaskUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Task entity.
func (tuo *TaskUpdateOne) Save(ctx context.Context) (*Task, error) {
	tuo.defaults()
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TaskUpdateOne) SaveX(ctx context.Context) *Task {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TaskUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TaskUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tuo *TaskUpdateOne) defaults() {
	if _, ok := tuo.mutation.UpdatedAt(); !ok {
		v := task.UpdateDefaultUpdatedAt()
		tuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tuo *TaskUpdateOne) check() error {
	if v, ok := tuo.mutation.Status(); ok {
		if err := task.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Task.status": %w`, err)}
		}
	}
	return nil
}

func (tuo *TaskUpdateOne) sqlSave(ctx context.Context) (_node *Task, err error) {
	if err := tuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(task.Table, task.Columns, sqlgraph.NewFieldSpec(task.FieldID, field.TypeString))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Task.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, task.FieldID)
		for _, f := range fields {
			if !task.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != task.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.Title(); ok {
		_spec.SetField(task.FieldTitle, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Description(); ok {
		_spec.SetField(task.FieldDescription, field.TypeString, value)
	}
	if tuo.mutation.DescriptionCleared() {
		_spec.ClearField(task.FieldDescription, field.TypeString)
	}
	if value, ok := tuo.mutation.DueDate(); ok {
		_spec.SetField(task.FieldDueDate, field.TypeTime, value)
	}
	if tuo.mutation.DueDateCleared() {
		_spec.ClearField(task.FieldDueDate, field.TypeTime)
	}
	if value, ok := tuo.mutation.AssigneeName(); ok {
		_spec.SetField(task.FieldAssigneeName, field.TypeString, value)
	}
	if tuo.mutation.AssigneeNameCleared() {
		_spec.ClearField(task.FieldAssigneeName, field.TypeString)
	}
	if value, ok := tuo.mutation.CreatedAt(); ok {
		_spec.SetField(task.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.UpdatedAt(); ok {
		_spec.SetField(task.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.CompletedAt(); ok {
		_spec.SetField(task.FieldCompletedAt, field.TypeTime, value)
	}
	if tuo.mutation.CompletedAtCleared() {
		_spec.ClearField(task.FieldCompletedAt, field.TypeTime)
	}
	if value, ok := tuo.mutation.Status(); ok {
		_spec.SetField(task.FieldStatus, field.TypeEnum, value)
	}
	if tuo.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   task.FilesTable,
			Columns: task.FilesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedFilesIDs(); len(nodes) > 0 && !tuo.mutation.FilesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   task.FilesTable,
			Columns: task.FilesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.FilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   task.FilesTable,
			Columns: task.FilesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Task{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{task.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
