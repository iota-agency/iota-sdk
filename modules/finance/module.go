package finance

import (
	"embed"
	"github.com/iota-uz/iota-sdk/components/icons"
	corepersistence "github.com/iota-uz/iota-sdk/modules/core/infrastructure/persistence"
	persistence2 "github.com/iota-uz/iota-sdk/modules/finance/infrastructure/persistence"
	"github.com/iota-uz/iota-sdk/modules/finance/permissions"
	"github.com/iota-uz/iota-sdk/modules/finance/presentation/controllers"
	"github.com/iota-uz/iota-sdk/modules/finance/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/spotlight"
)

//go:embed presentation/locales/*.json
var localeFiles embed.FS

//go:embed infrastructure/persistence/schema/finance-schema.sql
var migrationFiles embed.FS

func NewModule() application.Module {
	return &Module{}
}

type Module struct {
}

func (m *Module) Register(app application.Application) error {
	moneyAccountService := services.NewMoneyAccountService(
		persistence2.NewMoneyAccountRepository(),
		persistence2.NewTransactionRepository(),
		app.EventPublisher(),
	)
	currencyRepo := corepersistence.NewCurrencyRepository()
	transactionRepo := persistence2.NewTransactionRepository()
	categoryRepo := persistence2.NewExpenseCategoryRepository(currencyRepo)
	app.RegisterServices(
		services.NewPaymentService(
			persistence2.NewPaymentRepository(),
			app.EventPublisher(),
			moneyAccountService,
		),
		services.NewExpenseCategoryService(
			categoryRepo,
			app.EventPublisher(),
		),
		services.NewExpenseService(
			persistence2.NewExpenseRepository(categoryRepo, transactionRepo),
			app.EventPublisher(),
			moneyAccountService,
		),
		moneyAccountService,
		services.NewCounterpartyService(persistence2.NewCounterpartyRepository()),
	)

	app.RegisterControllers(
		controllers.NewExpensesController(app),
		controllers.NewMoneyAccountController(app),
		controllers.NewExpenseCategoriesController(app),
		controllers.NewPaymentsController(app),
		controllers.NewCounterpartiesController(app),
	)
	app.Spotlight().Register(
		spotlight.NewItem(nil, ExpenseCategoriesItem.Name, ExpenseCategoriesItem.Href),
		spotlight.NewItem(nil, PaymentsItem.Name, PaymentsItem.Href),
		spotlight.NewItem(nil, ExpensesItem.Name, ExpensesItem.Href),
		spotlight.NewItem(nil, AccountsItem.Name, AccountsItem.Href),
		spotlight.NewItem(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"Expenses.List.New",
			"/finance/expenses/new",
		),
		spotlight.NewItem(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"Accounts.List.New",
			"/finance/accounts/new",
		),
		spotlight.NewItem(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"Payments.List.New",
			"/finance/payments/new",
		),
		spotlight.NewItem(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"ExpenseCategories.List.New",
			"/finance/expense-categories/new",
		),
	)

	app.RegisterPermissions(permissions.Permissions...)
	app.RegisterLocaleFiles(&localeFiles)
	app.RegisterMigrationDirs(&migrationFiles)
	return nil
}

func (m *Module) Name() string {
	return "finance"
}
