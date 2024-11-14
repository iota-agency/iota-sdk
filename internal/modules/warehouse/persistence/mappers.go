package persistence

import (
	"errors"
	"github.com/iota-agency/iota-erp/internal/domain/aggregates/order"
	"github.com/iota-agency/iota-erp/internal/modules/warehouse/domain/entities/position"
	"github.com/iota-agency/iota-erp/internal/modules/warehouse/domain/entities/product"
	"github.com/iota-agency/iota-erp/internal/modules/warehouse/domain/entities/unit"
	"github.com/iota-agency/iota-erp/internal/modules/warehouse/persistence/models"
)

func toDBUnit(unit *unit.Unit) *models.WarehouseUnit {
	return &models.WarehouseUnit{
		ID:        unit.ID,
		Name:      unit.Name,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}
}

func toDomainUnit(dbUnit *models.WarehouseUnit) *unit.Unit {
	return &unit.Unit{
		ID:          dbUnit.ID,
		Name:        dbUnit.Name,
		Description: dbUnit.Description,
		CreatedAt:   dbUnit.CreatedAt,
		UpdatedAt:   dbUnit.UpdatedAt,
	}
}

func toDBOrder(data *order.Order) (*models.WarehouseOrder, []*models.OrderItem) {
	dbItems := make([]*models.OrderItem, 0, len(data.Items))
	for _, item := range data.Items {
		dbItems = append(
			dbItems, &models.OrderItem{
				ProductID: item.Product.ID,
				OrderID:   data.ID,
				CreatedAt: data.CreatedAt,
			},
		)
	}
	return &models.WarehouseOrder{
		ID:        data.ID,
		Status:    data.Status.String(),
		Type:      data.Type.String(),
		CreatedAt: data.CreatedAt,
	}, dbItems
}

func toDomainOrder(
	dbOrder *models.WarehouseOrder,
	dbItems []*models.OrderItem,
	dbProduct []*models.WarehouseProduct,
) (*order.Order, error) {
	items := make([]*order.Item, 0, len(dbItems))
	for _, item := range dbItems {
		var orderProduct *models.WarehouseProduct
		for _, p := range dbProduct {
			if p.ID == item.ProductID {
				orderProduct = p
				break
			}
		}
		if orderProduct == nil {
			return nil, errors.New("product not found")
		}
		p, err := toDomainProduct(orderProduct)
		if err != nil {
			return nil, err
		}
		items = append(
			items, &order.Item{
				Product:   p,
				CreatedAt: item.CreatedAt,
			},
		)
	}
	status, err := order.NewStatus(dbOrder.Status)
	if err != nil {
		return nil, err
	}
	typeEnum, err := order.NewType(dbOrder.Type)
	if err != nil {
		return nil, err
	}
	return &order.Order{
		ID:        dbOrder.ID,
		Status:    status,
		Type:      typeEnum,
		CreatedAt: dbOrder.CreatedAt,
		Items:     items,
	}, nil
}

func toDBProduct(entity *product.Product) *models.WarehouseProduct {
	return &models.WarehouseProduct{
		ID:         entity.ID,
		PositionID: entity.PositionID,
		Rfid:       entity.Rfid,
		Status:     string(entity.Status),
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}

func toDomainProduct(dbProduct *models.WarehouseProduct) (*product.Product, error) {
	status, err := product.NewStatus(dbProduct.Status)
	if err != nil {
		return nil, err
	}
	return &product.Product{
		ID:         dbProduct.ID,
		PositionID: dbProduct.PositionID,
		Rfid:       dbProduct.Rfid,
		Status:     status,
		CreatedAt:  dbProduct.CreatedAt,
		UpdatedAt:  dbProduct.UpdatedAt,
	}, nil
}

func toDomainPosition(dbPosition *models.WarehousePosition) (*position.Position, error) {
	return &position.Position{
		ID:        dbPosition.ID,
		Title:     dbPosition.Title,
		Barcode:   dbPosition.Barcode,
		UnitID:    dbPosition.UnitID,
		CreatedAt: dbPosition.CreatedAt,
		UpdatedAt: dbPosition.UpdatedAt,
	}, nil
}

func toDBPosition(entity *position.Position) *models.WarehousePosition {
	return &models.WarehousePosition{
		ID:        entity.ID,
		Title:     entity.Title,
		Barcode:   entity.Barcode,
		UnitID:    entity.UnitID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}