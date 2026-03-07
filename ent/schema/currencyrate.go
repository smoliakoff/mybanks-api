package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

type CurrencyRate struct {
	ent.Schema
}

func (CurrencyRate) Fields() []ent.Field {
	return []ent.Field{
		field.String("base").
			Default("GEL").
			Comment("Base currency ISO 4217 (usually GEL)"),
		field.String("currency").
			Comment("ISO 4217 currency code"),
		field.Float("buy").
			Optional().
			Nillable(),
		field.Float("sell").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(func() time.Time {
				return time.Now().UTC()
			}).
			Immutable(),
	}

}

func (CurrencyRate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bank", Bank.Type).
			Ref("currency_rates").
			Unique().
			Required(),
	}
}

func (CurrencyRate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}

func (CurrencyRate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("base", "currency"),
	}
}
