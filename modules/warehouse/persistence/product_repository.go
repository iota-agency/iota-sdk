package persistence

import (
	"context"
	"github.com/iota-agency/iota-sdk/modules/warehouse/domain/aggregates/product"
	"github.com/iota-agency/iota-sdk/modules/warehouse/persistence/models"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/graphql/helpers"
	"github.com/iota-agency/iota-sdk/pkg/mapping"
	"gorm.io/gorm"
)

type GormProductRepository struct{}

func NewProductRepository() product.Repository {
	return &GormProductRepository{}
}

func (g *GormProductRepository) tx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	return tx.Preload("Position"), nil
}

func (g *GormProductRepository) GetPaginated(
	ctx context.Context, limit, offset int,
	sortBy []string,
) ([]*product.Product, error) {
	tx, err := g.tx(ctx)
	if err != nil {
		return nil, err
	}
	q := tx.Limit(limit).Offset(offset)
	q, err = helpers.ApplySort(q, sortBy)
	if err != nil {
		return nil, err
	}
	var entities []*models.WarehouseProduct
	if err := q.Find(&entities).Error; err != nil {
		return nil, err
	}
	products := make([]*product.Product, len(entities))
	for i, entity := range entities {
		p, err := toDomainProduct(entity)
		if err != nil {
			return nil, err
		}
		products[i] = p
	}
	return products, nil
}

func (g *GormProductRepository) Count(ctx context.Context) (int64, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return 0, composables.ErrNoTx
	}
	var count int64
	if err := tx.Model(&models.WarehouseProduct{}).Count(&count).Error; err != nil { //nolint:exhaustruct
		return 0, err
	}
	return count, nil
}

func (g *GormProductRepository) GetAll(ctx context.Context) ([]*product.Product, error) {
	tx, err := g.tx(ctx)
	if err != nil {
		return nil, err
	}
	var entities []*models.WarehouseProduct
	if err := tx.Find(&entities).Error; err != nil {
		return nil, err
	}
	products := make([]*product.Product, len(entities))
	for i, entity := range entities {
		p, err := toDomainProduct(entity)
		if err != nil {
			return nil, err
		}
		products[i] = p
	}
	return products, nil
}

func (g *GormProductRepository) GetByID(ctx context.Context, id uint) (*product.Product, error) {
	tx, err := g.tx(ctx)
	if err != nil {
		return nil, err
	}
	var entity models.WarehouseProduct
	if err := tx.Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}
	return toDomainProduct(&entity)
}

func (g *GormProductRepository) GetByRfid(ctx context.Context, rfid string) (*product.Product, error) {
	tx, err := g.tx(ctx)
	if err != nil {
		return nil, err
	}
	var entity models.WarehouseProduct
	if err := tx.Where("rfid = ?", rfid).First(&entity).Error; err != nil {
		return nil, err
	}
	return toDomainProduct(&entity)
}

func (g *GormProductRepository) Create(ctx context.Context, data *product.Product) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	dbRow, err := toDBProduct(data)
	if err != nil {
		return err
	}
	return tx.Create(dbRow).Error
}

func (g *GormProductRepository) BulkCreate(ctx context.Context, data []*product.Product) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	dbRows, err := mapping.MapDbModels(data, toDBProduct)
	if err != nil {
		return err
	}
	maxParams := 1000
	for i := 0; i < len(dbRows); i += maxParams {
		end := i + maxParams
		if end > len(dbRows) {
			end = len(dbRows)
		}
		if err := tx.Create(dbRows[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (g *GormProductRepository) CreateOrUpdate(ctx context.Context, data *product.Product) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	dbRow, err := toDBProduct(data)
	if err != nil {
		return err
	}
	return tx.Save(dbRow).Error
}

func (g *GormProductRepository) Update(ctx context.Context, data *product.Product) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	dbRow, err := toDBProduct(data)
	if err != nil {
		return err
	}
	return tx.Save(dbRow).Error
}

func (g *GormProductRepository) Delete(ctx context.Context, id uint) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	return tx.Where("id = ?", id).Delete(&models.WarehouseProduct{}).Error //nolint:exhaustruct
}
