package mixin

import (
	"entgo.io/ent"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/mixin"
	"github.com/NpoolPlatform/stock-manager/pkg/db/ent/privacy"
	"github.com/NpoolPlatform/stock-manager/pkg/db/rule"
)

func (TimeMixin) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (TimeMixin) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			rule.FilterTimeRule(),
		},
	}
}
