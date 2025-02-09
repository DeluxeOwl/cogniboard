package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("title"),
		field.String("description").
			Optional().
			Nillable(),
		field.Time("due_date").
			Optional().
			Nillable(),
		field.String("assignee_name").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("completed_at").
			Optional().
			Nillable(),
		field.Enum("status").
			Values(
				"pending",
				"in_progress",
				"in_review",
				"completed",
			).
			Default("pending"),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
	}
}
