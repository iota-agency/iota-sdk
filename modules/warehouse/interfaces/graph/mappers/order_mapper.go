package mappers

import (
	"github.com/iota-agency/iota-sdk/modules/warehouse/domain/aggregates/order"
	"github.com/iota-agency/iota-sdk/modules/warehouse/domain/aggregates/product"
	model "github.com/iota-agency/iota-sdk/modules/warehouse/interfaces/graph/gqlmodels"
	"github.com/iota-agency/iota-sdk/pkg/mapping"
)

func OrderItemsToGraphModel(item order.Item) *model.OrderItem {
	pos := item.Position()
	return &model.OrderItem{
		Position: PositionToGraphModel(&pos),
		Products: mapping.MapViewModels(item.Products(), func(p *product.Product) *model.Product {
			return ProductToGraphModel(p)
		}),
	}
}

func OrderToGraphModel(o order.Order) *model.Order {
	return &model.Order{
		ID:        int64(o.ID()),
		Type:      string(o.Type()),
		Status:    string(o.Status()),
		Items:     mapping.MapViewModels(o.Items(), OrderItemsToGraphModel),
		CreatedAt: o.CreatedAt(),
	}
}
