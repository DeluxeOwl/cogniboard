package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("name"),
		field.Int64("size"),
		field.String("mime_type"),
		field.Time("uploaded_at").
			Default(time.Now),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).
			Ref("files"),
	}
}
