package graph

import (
	"context"
	"mybanks-api/ent"
	"mybanks-api/ent/bank"
	"mybanks-api/ent/banktranslation"
	"mybanks-api/ent/currencyrate"
	"mybanks-api/graph/model"
)

// Resolver is the root dependency container for GraphQL resolvers.
// gqlgen will use it to satisfy generated resolver interfaces.
type Resolver struct {
	Client *ent.Client
}

func (r *Resolver) Mutation() MutationResolver {
	return r
}

func (r *Resolver) ImportCurrencyRates(
	ctx context.Context,
	bankID int,
	rates []*model.ImportCurrencyRateInput,
	replaceExisting *bool,
) ([]*ent.CurrencyRate, error) {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	rollback := func(err error) ([]*ent.CurrencyRate, error) {
		_ = tx.Rollback()
		return nil, err
	}

	shouldReplace := true
	if replaceExisting != nil {
		shouldReplace = *replaceExisting
	}

	if shouldReplace {
		if _, err := tx.CurrencyRate.
			Delete().
			Where(currencyrate.HasBankWith(bank.IDEQ(bankID))).
			Exec(ctx); err != nil {
			return rollback(err)
		}
	}

	created := make([]*ent.CurrencyRate, 0, len(rates))

	for _, rate := range rates {
		builder := tx.CurrencyRate.
			Create().
			SetBankID(bankID).
			SetBase(rate.Base).
			SetCurrency(rate.Currency)

		if rate.Buy != nil {
			builder.SetBuy(*rate.Buy)
		}
		if rate.Sell != nil {
			builder.SetSell(*rate.Sell)
		}
		if rate.CreatedAt != nil {
			builder.SetCreatedAt(*rate.CreatedAt)
		}

		row, err := builder.Save(ctx)
		if err != nil {
			return rollback(err)
		}
		created = append(created, row)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return created, nil
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
