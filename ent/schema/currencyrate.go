package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CurrencyRate struct {
	ent.Schema
}

func (CurrencyRate) Fields() []ent.Field {
	return []ent.Field{
		field.String("currency").
			Comment("ISO 4217 currency code"),
		field.Float("buy").
			Optional().
			Nillable(),
		field.Float("sell").
			Optional().
			Nillable(),
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
