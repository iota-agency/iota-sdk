package expense

import (
	category "github.com/iota-agency/iota-erp/internal/domain/expense_category"
	"github.com/iota-agency/iota-erp/internal/interfaces/graph/gqlmodels"
	"time"
)

type Expense struct {
	Id         int64 `gorm:"primaryKey"`
	Amount     float64
	CategoryId int64
	Category   *category.ExpenseCategory `gorm:"foreignKey:CategoryId"`
	Date       time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (e *Expense) category2Graph() *model.ExpenseCategory {
	if e.Category == nil {
		return nil
	}
	return e.Category.ToGraph()
}

func (e *Expense) ToGraph() *model.Expense {
	return &model.Expense{
		ID:         e.Id,
		Amount:     e.Amount,
		CategoryID: e.CategoryId,
		Category:   e.category2Graph(),
		Date:       e.Date,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}