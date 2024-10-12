package mappers

import (
	"fmt"
	category "github.com/iota-agency/iota-erp/internal/domain/aggregates/expense_category"
	moneyAccount "github.com/iota-agency/iota-erp/internal/domain/aggregates/money_account"
	"github.com/iota-agency/iota-erp/internal/domain/aggregates/payment"
	stage "github.com/iota-agency/iota-erp/internal/domain/entities/project_stages"
	"github.com/iota-agency/iota-erp/internal/presentation/viewmodels"
	"strconv"
	"time"
)

func ExpenseCategoryToViewModel(entity *category.ExpenseCategory) *viewmodels.ExpenseCategory {
	return &viewmodels.ExpenseCategory{
		ID:                 strconv.FormatUint(uint64(entity.ID), 10),
		Name:               entity.Name,
		Amount:             fmt.Sprintf("%.2f", entity.Amount),
		AmountWithCurrency: fmt.Sprintf("%.2f %s", entity.Amount, entity.Currency.Symbol.String()),
		Description:        entity.Description,
		UpdatedAt:          entity.UpdatedAt.Format(time.RFC3339),
		CreatedAt:          entity.CreatedAt.Format(time.RFC3339),
	}
}

func MoneyAccountToViewModel(entity *moneyAccount.Account) *viewmodels.MoneyAccount {
	return &viewmodels.MoneyAccount{
		ID:                  strconv.FormatUint(uint64(entity.ID), 10),
		Name:                entity.Name,
		AccountNumber:       entity.AccountNumber,
		Balance:             fmt.Sprintf("%.2f", entity.Balance),
		BalanceWithCurrency: fmt.Sprintf("%.2f %s", entity.Balance, entity.Currency.Symbol.String()),
		CurrencyCode:        entity.Currency.Code.String(),
		CurrencySymbol:      entity.Currency.Symbol.String(),
		Description:         entity.Description,
		UpdatedAt:           entity.UpdatedAt.Format(time.RFC3339),
		CreatedAt:           entity.CreatedAt.Format(time.RFC3339),
	}
}

func ProjectStageToViewModel(entity *stage.ProjectStage) *viewmodels.ProjectStage {
	return &viewmodels.ProjectStage{
		ID:        strconv.FormatUint(uint64(entity.ID), 10),
		Name:      entity.Name,
		ProjectID: strconv.FormatUint(uint64(entity.ProjectID), 10),
		Margin:    fmt.Sprintf("%.2f", entity.Margin),
		Risks:     fmt.Sprintf("%.2f", entity.Risks),
		UpdatedAt: entity.UpdatedAt.Format(time.RFC3339),
		CreatedAt: entity.CreatedAt.Format(time.RFC3339),
	}
}

func PaymentToViewModel(entity *payment.Payment) *viewmodels.Payment {
	currency := entity.Account.Currency
	return &viewmodels.Payment{
		ID:                 strconv.FormatUint(uint64(entity.ID), 10),
		Amount:             fmt.Sprintf("%.2f", entity.Amount),
		AmountWithCurrency: fmt.Sprintf("%.2f %s", entity.Amount, currency.Symbol.String()),
		AccountID:          strconv.FormatUint(uint64(entity.Account.ID), 10),
		TransactionID:      strconv.FormatUint(uint64(entity.TransactionID), 10),
		TransactionDate:    entity.TransactionDate.Format(time.RFC3339),
		AccountingPeriod:   entity.AccountingPeriod.Format(time.RFC3339),
		Comment:            entity.Comment,
		CreatedAt:          entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          entity.UpdatedAt.Format(time.RFC3339),
	}
}