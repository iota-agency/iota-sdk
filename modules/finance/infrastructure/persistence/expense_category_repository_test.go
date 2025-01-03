package persistence_test

//import (
//	category "github.com/iota-uz/iota-sdk/modules/finance/domain/aggregates/expense_category"
//	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/currency"
//	"github.com/iota-uz/iota-sdk/pkg/infrastructure/persistence"
//	"github.com/iota-uz/iota-sdk/pkg/testutils"
//	"testing"
//	"time"
//)
//
//func TestGormExpenseCategoryRepository_CRUD(t *testing.T) {
//	ctx := testutils.GetTestContext()
//	currencyRepository := persistence.NewCurrencyRepository()
//	categoryRepository := persistence.NewExpenseCategoryRepository()
//
//	if err := currencyRepository.Create(ctx.Context, &currency.USD); err != nil {
//		t.Fatal(err)
//	}
//	if err := categoryRepository.Create(
//		ctx.Context, &category.ExpenseCategory{
//			ID:          1,
//			Path:        "test",
//			Amount:      100,
//			Currency:    currency.USD,
//			Description: "test",
//			CreatedAt:   time.Now(),
//			UpdatedAt:   time.Now(),
//		},
//	); err != nil {
//		t.Fatal(err)
//	}
//
//	t.Run("Count", func(t *testing.T) {
//		count, err := categoryRepository.Count(ctx.Context)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if count != 1 {
//			t.Errorf("expected 1, got %d", count)
//		}
//	})
//
//	t.Run("GetPaginated", func(t *testing.T) {
//		categories, err := categoryRepository.GetPaginated(ctx.Context, 1, 0, []string{})
//		if err != nil {
//			t.Fatal(err)
//		}
//		if len(categories) != 1 {
//			t.Errorf("expected 1, got %d", len(categories))
//		}
//		if categories[0].Amount != 100 {
//			t.Errorf("expected 100, got %f", categories[0].Amount)
//		}
//	})
//
//	t.Run("GetAll", func(t *testing.T) {
//		categories, err := categoryRepository.GetAll(ctx.Context)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if len(categories) != 1 {
//			t.Errorf("expected 1, got %d", len(categories))
//		}
//		if categories[0].Amount != 100 {
//			t.Errorf("expected 100, got %f", categories[0].Amount)
//		}
//	})
//
//	t.Run("GetByID", func(t *testing.T) {
//		categoryEntity, err := categoryRepository.GetByID(ctx.Context, 1)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if categoryEntity.Amount != 100 {
//			t.Errorf("expected 100, got %f", categoryEntity.Amount)
//		}
//		if categoryEntity.Currency.Code.Get() != currency.UsdCode {
//			t.Errorf("expected %s, got %s", currency.UsdCode, categoryEntity.Currency.Code.Get())
//		}
//	})
//
//	t.Run("Update", func(t *testing.T) {
//		if err := categoryRepository.Update(
//			ctx.Context, &category.ExpenseCategory{
//				ID:     1,
//				Amount: 200,
//			},
//		); err != nil {
//			t.Fatal(err)
//		}
//		categoryEntity, err := categoryRepository.GetByID(ctx.Context, 1)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if categoryEntity.Amount != 200 {
//			t.Errorf("expected 200, got %f", categoryEntity.Amount)
//		}
//		if categoryEntity.Currency.Code.Get() != currency.UsdCode {
//			t.Errorf("expected %s, got %s", currency.UsdCode, categoryEntity.Currency.Code.Get())
//		}
//	})
//}
