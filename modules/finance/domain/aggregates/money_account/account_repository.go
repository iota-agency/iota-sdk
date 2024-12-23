package moneyaccount

import (
	"context"
)

type DateRange struct {
	From string
	To   string
}

type FindParams struct {
	Limit     int
	Offset    int
	SortBy    []string
	Query     string
	Field     string
	CreatedAt DateRange
}

type Repository interface {
	Count(ctx context.Context) (uint, error)
	GetAll(ctx context.Context) ([]*Account, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]*Account, error)
	GetByID(ctx context.Context, id uint) (*Account, error)
	RecalculateBalance(ctx context.Context, id uint) error
	Create(ctx context.Context, payment *Account) error
	Update(ctx context.Context, payment *Account) error
	Delete(ctx context.Context, id uint) error
}
