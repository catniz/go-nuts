package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
)

// SystemConfig holds the schema definition for the SystemConfig entity.
type SystemConfig struct {
	ent.Schema
}

// Fields of the SystemConfig.
func (SystemConfig) Fields() []ent.Field {
	return []ent.Field{
		// Fields go here.
	}
}

// Edges of the SystemConfig.
func (SystemConfig) Edges() []ent.Edge {
	return []ent.Edge{
		// Edges go here.
	}
}

// Mixin of the SystemConfig.
func (SystemConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IdMixin{},
		TimeMixin{},
	}
}

// Annotations of the SystemConfig.
func (SystemConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "system_config", Schema: "bullish"}}
}
