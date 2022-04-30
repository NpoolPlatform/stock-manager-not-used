package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/NpoolPlatform/stock-manager/pkg/db/mixin"
	"github.com/google/uuid"
)

// Stock holds the schema definition for the Stock entity.
type Stock struct {
	ent.Schema
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
		field.Uint32("total"),
		field.Uint32("locked"),
		field.Uint32("in_service"),
		field.Uint32("sold"),
	}
}

func (Stock) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("good_id").Unique(),
	}
}
