package graph

import (
	"context"
	"mybanks-api/ent/banktranslation"

	"mybanks-api/ent"
)

// Resolver is the root dependency container for GraphQL resolvers.
// gqlgen will use it to satisfy generated resolver interfaces.
type Resolver struct {
	Client *ent.Client
}

func (r *Resolver) Translation(ctx context.Context, obj *ent.Bank, locale string) (*ent.BankTranslation, error) {
	if obj == nil {
		return nil, nil
	}

	tr, err := r.Client.BankTranslation.
		Query().
		Where(
			banktranslation.BankIDEQ(obj.ID),
			banktranslation.LocaleEQ(locale),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return tr, nil
}

// Query returns implementation of the generated QueryResolver interface.
func (r *Resolver) Query() QueryResolver { return r }

// --------------------
// Query resolvers
// --------------------

func (r *Resolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.Client.Noder(ctx, id)
}

func (r *Resolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.Client.Noders(ctx, ids)
}

func (r *Resolver) Banks(
	ctx context.Context,
	after *ent.Cursor,
	first *int,
	before *ent.Cursor,
	last *int,
	where *ent.BankWhereInput,
) (*ent.BankConnection, error) {
	q := r.Client.Bank.Query()

	if where != nil {
		var err error
		q, err = where.Filter(q)
		if err != nil {
			return nil, err
		}
	}

	return q.Paginate(ctx, after, first, before, last)
}

func (r *Resolver) CurrencyRates(ctx context.Context) ([]*ent.CurrencyRate, error) {
	return r.Client.CurrencyRate.Query().All(ctx)
}

func (r *Resolver) Offers(ctx context.Context) ([]*ent.Offer, error) {
	return r.Client.Offer.Query().All(ctx)
}

// Bank returns implementation of the generated BankResolver interface.
func (r *Resolver) Bank() BankResolver { return r }
