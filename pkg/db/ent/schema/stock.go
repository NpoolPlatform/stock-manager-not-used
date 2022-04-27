package schema

import "entgo.io/ent"

// Stock holds the schema definition for the Stock entity.
type Stock struct {
	ent.Schema
}

// Fields of the Stock.
func (Stock) Fields() []ent.Field {
	return nil
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
	return nil
}
