package persistence

import (
	"context"
	"github.com/iota-agency/iota-sdk/modules/warehouse/domain/aggregates/order"
	"github.com/iota-agency/iota-sdk/modules/warehouse/persistence/models"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/graphql/helpers"
	"github.com/iota-agency/iota-sdk/pkg/mapping"
	"gorm.io/gorm"
)

type GormOrderRepository struct{}

func NewOrderRepository() order.Repository {
	return &GormOrderRepository{}
}

func (g *GormOrderRepository) tx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	return tx, nil
}

func (g *GormOrderRepository) GetPaginated(ctx context.Context, params *order.FindParams) ([]*order.Order, error) {
	tx, err := g.tx(ctx)
	if err != nil {
		return nil, err
	}
	q := tx.Limit(params.Limit).Offset(params.Offset)
	q, err = helpers.ApplySort(q, params.SortBy)
	if err != nil {
		return nil, err
	}
	var rows []*models.WarehouseOrder
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	for i, row := range rows {
		products, err := g.getProducts(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		rows[i].Products = products
	}
	return mapping.MapDbModels(rows, toDomainOrder)
}

func (g *GormOrderRepository) getProducts(ctx context.Context, id uint) ([]*models.WarehouseProduct, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var entities []*models.WarehouseOrderItem
	if err := tx.Where("warehouse_order_id = ?", id).Find(&entities).Error; err != nil {
		return nil, err
	}
	var rows []*models.WarehouseProduct
	for _, entity := range entities {
		var product models.WarehouseProduct
		if err := tx.Where("id = ?", entity.ProductID).First(&product).Error; err != nil {
			return nil, err
		}
		rows = append(rows, &product)
	}
	return rows, nil
}

func (g *GormOrderRepository) Count(ctx context.Context) (int64, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return 0, composables.ErrNoTx
	}
	var count int64
	if err := tx.Model(&models.WarehouseOrder{}).Count(&count).Error; err != nil { //nolint:exhaustruct
		return 0, err
	}
	return count, nil
}

func (g *GormOrderRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	// TODO: proper implementation
	var rows []*models.WarehouseOrder
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	for i, row := range rows {
		products, err := g.getProducts(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		rows[i].Products = products
	}
	return mapping.MapDbModels(rows, toDomainOrder)
}

func (g *GormOrderRepository) GetByID(ctx context.Context, id uint) (*order.Order, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var row models.WarehouseOrder
	if err := tx.Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	products, err := g.getProducts(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	row.Products = products
	return toDomainOrder(&row)
}

func (g *GormOrderRepository) Create(ctx context.Context, data *order.Order) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	or, orderItems := toDBOrder(data)
	if err := tx.Create(or).Error; err != nil {
		return err
	}
	for _, item := range orderItems {
		item.WarehouseOrderID = or.ID
	}
	return tx.Create(orderItems).Error
}

func (g *GormOrderRepository) Update(ctx context.Context, data *order.Order) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	or, orderItems := toDBOrder(data)
	if err := tx.Updates(or).Error; err != nil {
		return err
	}
	if err := tx.Where("warehouse_order_id = ?", or.ID).Delete(&models.WarehouseOrderItem{}).Error; err != nil { //nolint:exhaustruct
		return err
	}
	if len(orderItems) == 0 {
		return nil
	}
	return tx.Create(orderItems).Error
}

func (g *GormOrderRepository) Delete(ctx context.Context, id uint) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	if err := tx.Where("id = ?", id).Delete(&models.WarehouseOrder{}).Error; err != nil { //nolint:exhaustruct
		return err
	}
	return nil
}
