package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.String("name").Default("unknown"),
		field.String("email").Unique(),
		field.String("password").Default("unknown"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
