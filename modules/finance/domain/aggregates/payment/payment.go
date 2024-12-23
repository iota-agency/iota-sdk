package payment

import (
	moneyAccount "github.com/iota-agency/iota-sdk/modules/finance/domain/aggregates/money_account"
	"github.com/iota-agency/iota-sdk/pkg/domain/aggregates/user"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Payment struct {
	ID               uint
	StageID          uint
	Amount           float64
	TransactionID    uint
	Account          moneyAccount.Account
	TransactionDate  time.Time
	AccountingPeriod time.Time
	Comment          string
	User             *user.User
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateDTO struct {
	Amount           float64   `validate:"required,gt=0"`
	AccountID        uint      `validate:"required"`
	TransactionDate  time.Time `validate:"required"`
	AccountingPeriod time.Time `validate:"required"`
	Comment          string
	UserID           uint `validate:"required"`
	StageID          uint `validate:"required"`
}

type UpdateDTO struct {
	Amount           float64 `validate:"gt=0"`
	AccountID        uint
	TransactionDate  time.Time
	AccountingPeriod time.Time
	Comment          string
	UserID           uint
	StageID          uint
}

func (p *CreateDTO) Ok(l ut.Translator) (map[string]string, bool) {
	errors := map[string]string{}
	err := validate.Struct(p)
	if err == nil {
		return errors, true
	}
	for _, _err := range err.(validator.ValidationErrors) {
		errors[_err.Field()] = _err.Translate(l)
	}
	return errors, len(errors) == 0
}

func (p *CreateDTO) ToEntity() *Payment {
	return &Payment{
		ID:               0,
		TransactionID:    0,
		StageID:          p.StageID,
		Amount:           p.Amount,
		Account:          moneyAccount.Account{ID: p.AccountID}, //nolint:exhaustruct
		TransactionDate:  p.TransactionDate,
		AccountingPeriod: p.AccountingPeriod,
		User:             &user.User{ID: p.UserID}, //nolint:exhaustruct
		Comment:          p.Comment,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (p *Payment) Ok(l ut.Translator) (map[string]string, bool) {
	errors := map[string]string{}
	errs := validate.Struct(p)
	if errs == nil {
		return errors, true
	}
	for _, err := range errs.(validator.ValidationErrors) {
		errors[err.Field()] = err.Translate(l)
	}
	return errors, len(errors) == 0
}

func (p *UpdateDTO) Ok(l ut.Translator) (map[string]string, bool) {
	errors := map[string]string{}
	errs := validate.Struct(p)
	if errs == nil {
		return errors, true
	}
	for _, err := range errs.(validator.ValidationErrors) {
		errors[err.Field()] = err.Translate(l)
	}
	return errors, len(errors) == 0
}

func (p *UpdateDTO) ToEntity(id uint) *Payment {
	return &Payment{
		ID:               id,
		StageID:          p.StageID,
		Amount:           p.Amount,
		Account:          moneyAccount.Account{ID: p.AccountID}, //nolint:exhaustruct
		TransactionDate:  p.TransactionDate,
		TransactionID:    0,
		AccountingPeriod: p.AccountingPeriod,
		Comment:          p.Comment,
		User:             &user.User{ID: p.UserID}, //nolint:exhaustruct
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
