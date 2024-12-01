package seed

import (
	"context"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/currency"
	"github.com/iota-agency/iota-sdk/pkg/infrastructure/persistence"
)

func CreateCurrencies(ctx context.Context, app application.Application) error {
	currencyRepository := persistence.NewCurrencyRepository()
	for _, c := range currency.Currencies {
		if err := currencyRepository.CreateOrUpdate(ctx, &c); err != nil {
			return err
		}
	}
	return nil
}
