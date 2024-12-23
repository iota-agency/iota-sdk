package persistence

import (
	"context"
	"github.com/iota-agency/iota-sdk/modules/finance/domain/entities/transaction"
	"github.com/iota-agency/iota-sdk/modules/finance/persistence/models"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/graphql/helpers"
)

type GormTransactionRepository struct{}

func NewTransactionRepository() transaction.Repository {
	return &GormTransactionRepository{}
}

func (g *GormTransactionRepository) GetPaginated(
	ctx context.Context,
	limit, offset int,
	sortBy []string,
) ([]*transaction.Transaction, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	q := tx.Limit(limit).Offset(offset)
	q, err := helpers.ApplySort(q, sortBy)
	if err != nil {
		return nil, err
	}
	var entities []*transaction.Transaction
	if err := q.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (g *GormTransactionRepository) Count(ctx context.Context) (int64, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return 0, composables.ErrNoTx
	}
	var count int64
	if err := tx.Model(&transaction.Transaction{}).Count(&count).Error; err != nil { //nolint:exhaustruct
		return 0, err
	}
	return count, nil
}

func (g *GormTransactionRepository) GetAll(ctx context.Context) ([]*transaction.Transaction, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var entities []*transaction.Transaction
	if err := tx.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (g *GormTransactionRepository) GetByID(ctx context.Context, id int64) (*transaction.Transaction, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var entity models.Transaction
	if err := tx.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return toDomainTransaction(&entity)
}

func (g *GormTransactionRepository) Create(ctx context.Context, data *transaction.Transaction) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	entity := toDBTransaction(data)
	return tx.Create(entity).Error
}

func (g *GormTransactionRepository) Update(ctx context.Context, data *transaction.Transaction) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	entity := toDBTransaction(data)
	return tx.Save(entity).Error
}

func (g *GormTransactionRepository) Delete(ctx context.Context, id int64) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	return tx.Delete(&transaction.Transaction{}, id).Error //nolint:exhaustruct
}
