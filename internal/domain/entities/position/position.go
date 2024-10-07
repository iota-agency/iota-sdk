package position

import (
	"time"

	model "github.com/iota-agency/iota-erp/internal/interfaces/graph/gqlmodels"
)

type Position struct {
	Id          int64
	Name        string
	Description *string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (p *Position) ToGraph() *model.Position {
	return &model.Position{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   *p.CreatedAt,
		UpdatedAt:   *p.UpdatedAt,
	}
}
