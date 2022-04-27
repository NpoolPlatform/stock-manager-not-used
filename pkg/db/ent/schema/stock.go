package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"

	"github.com/NpoolPlatform/stock-manager/pkg/db/mixin"
)

// Stock holds the schema definition for the Stock entity.
type Stock struct {
	ent.Schema
	mixin.TimeMixin
}

func (Stock) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

// Fields of the Stock.
func (Stock) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("good_id", uuid.UUID{}),
		field.Int32("total"),
		field.Int32("in_service"),
		field.Int32("sold"),
	}
}

func (Stock) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("good_id").Unique(),
	}
}
