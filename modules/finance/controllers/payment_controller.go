package controllers

import (
	"github.com/a-h/templ"
	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/iota-agency/iota-sdk/modules/finance/domain/aggregates/payment"
	"github.com/iota-agency/iota-sdk/modules/finance/services"
	"github.com/iota-agency/iota-sdk/modules/finance/templates/pages/payments"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/mapping"
	"github.com/iota-agency/iota-sdk/pkg/shared"
	"github.com/iota-agency/iota-sdk/pkg/shared/middleware"
	"github.com/iota-agency/iota-sdk/pkg/types"
	"net/http"

	"github.com/iota-agency/iota-sdk/modules/finance/mappers"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/presentation/viewmodels"
)

type PaymentsController struct {
	app                 application.Application
	paymentService      *services.PaymentService
	moneyAccountService *services.MoneyAccountService
	basePath            string
}

func NewPaymentsController(app application.Application) application.Controller {
	return &PaymentsController{
		app:                 app,
		paymentService:      app.Service(services.PaymentService{}).(*services.PaymentService),
		moneyAccountService: app.Service(services.MoneyAccountService{}).(*services.MoneyAccountService),
		basePath:            "/finance/payments",
	}
}

func (c *PaymentsController) Register(r *mux.Router) {
	router := r.PathPrefix(c.basePath).Subrouter()
	router.Use(middleware.RequireAuthorization())
	router.HandleFunc("", c.Payments).Methods(http.MethodGet)
	router.HandleFunc("", c.CreatePayment).Methods(http.MethodPost)
	router.HandleFunc("/new", c.GetNew).Methods(http.MethodGet)
	router.HandleFunc("/{id:[0-9]+}", c.GetEdit).Methods(http.MethodGet)
	router.HandleFunc("/{id:[0-9]+}", c.PostEdit).Methods(http.MethodPost)
	router.HandleFunc("/{id:[0-9]+}", c.DeletePayment).Methods(http.MethodDelete)
}

func (c *PaymentsController) viewModelPayments(r *http.Request) ([]*viewmodels.Payment, error) {
	ps, err := c.paymentService.GetAll(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving payments")
	}
	vms := make([]*viewmodels.Payment, len(ps))
	for i, p := range ps {
		vms[i] = mappers.PaymentToViewModel(p)
	}
	return vms, nil
}

func (c *PaymentsController) viewModelPayment(r *http.Request) (*viewmodels.Payment, error) {
	id, err := shared.ParseID(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing id")
	}
	p, err := c.paymentService.GetByID(r.Context(), id)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving payment")
	}
	return mappers.PaymentToViewModel(p), nil
}

func (c *PaymentsController) viewModelAccounts(r *http.Request) ([]*viewmodels.MoneyAccount, error) {
	accounts, err := c.moneyAccountService.GetAll(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving accounts")
	}
	return mapping.MapViewModels(accounts, mappers.MoneyAccountToViewModel), nil
}

func (c *PaymentsController) Payments(w http.ResponseWriter, r *http.Request) {
	pageCtx, err := composables.UsePageCtx(
		r,
		types.NewPageData("Payments.Meta.List.Title", ""),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	paymentViewModels, err := c.viewModelPayments(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isHxRequest := len(r.Header.Get("Hx-Request")) > 0
	props := &payments.IndexPageProps{
		PageContext: pageCtx,
		Payments:    paymentViewModels,
	}
	if isHxRequest {
		templ.Handler(payments.PaymentsTable(props), templ.WithStreaming()).ServeHTTP(w, r)
	} else {
		templ.Handler(payments.Index(props), templ.WithStreaming()).ServeHTTP(w, r)
	}
}

func (c *PaymentsController) GetEdit(w http.ResponseWriter, r *http.Request) {
	pageCtx, err := composables.UsePageCtx(
		r,
		types.NewPageData("Payments.Meta.Edit.Title", ""),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	paymentViewModel, err := c.viewModelPayment(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accounts, err := c.viewModelAccounts(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := &payments.EditPageProps{
		PageContext: pageCtx,
		Payment:     paymentViewModel,
		Accounts:    accounts,
		Errors:      make(map[string]string),
	}
	templ.Handler(payments.Edit(props), templ.WithStreaming()).ServeHTTP(w, r)
}

func (c *PaymentsController) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id, err := shared.ParseID(r)
	if err != nil {
		http.Error(w, "Error parsing id", http.StatusInternalServerError)
		return
	}

	if _, err := c.paymentService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shared.Redirect(w, r, c.basePath)
}

func (c *PaymentsController) PostEdit(w http.ResponseWriter, r *http.Request) {
	id, err := shared.ParseID(r)
	if err != nil {
		http.Error(w, "Error parsing id", http.StatusInternalServerError)
		return
	}
	action := shared.FormAction(r.FormValue("_action"))
	if !action.IsValid() {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}
	r.Form.Del("_action")

	switch action {
	case shared.FormActionDelete:
		_, err = c.paymentService.Delete(r.Context(), id)
	case shared.FormActionSave:
		dto := payment.UpdateDTO{} //nolint:exhaustruct
		if err := shared.Decoder.Decode(&dto, r.Form); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var pageCtx *types.PageContext
		pageCtx, err = composables.UsePageCtx(
			r,
			types.NewPageData("Payments.Meta.Edit.Title", ""),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if errorsMap, ok := dto.Ok(pageCtx.UniTranslator); !ok {
			paymentViewModel, err := c.viewModelPayment(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			accounts, err := c.viewModelAccounts(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			props := &payments.EditPageProps{
				PageContext: pageCtx,
				Payment:     paymentViewModel,
				Accounts:    accounts,
				Errors:      errorsMap,
			}
			templ.Handler(payments.EditForm(props), templ.WithStreaming()).ServeHTTP(w, r)
			return
		}
		err = c.paymentService.Update(r.Context(), id, &dto)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shared.Redirect(w, r, c.basePath)
}

func (c *PaymentsController) GetNew(w http.ResponseWriter, r *http.Request) {
	pageCtx, err := composables.UsePageCtx(
		r, types.NewPageData("Payments.Meta.New.Title", ""),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accounts, err := c.viewModelAccounts(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := &payments.CreatePageProps{
		PageContext: pageCtx,
		Payment:     &viewmodels.Payment{}, //nolint:exhaustruct
		Accounts:    accounts,
		Errors:      make(map[string]string),
	}
	templ.Handler(payments.New(props), templ.WithStreaming()).ServeHTTP(w, r)
}

func (c *PaymentsController) CreatePayment(w http.ResponseWriter, r *http.Request) {
	dto := payment.CreateDTO{} //nolint:exhaustruct
	if err := shared.Decoder.Decode(&dto, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageCtx, err := composables.UsePageCtx(
		r, types.NewPageData("Payments.Meta.New.Title", ""),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errorsMap, ok := dto.Ok(pageCtx.UniTranslator); !ok {
		accounts, err := c.viewModelAccounts(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		props := &payments.CreatePageProps{
			PageContext: pageCtx,
			Payment:     mappers.PaymentToViewModel(dto.ToEntity()),
			Accounts:    accounts,
			Errors:      errorsMap,
		}
		templ.Handler(
			payments.CreateForm(props),
			templ.WithStreaming(),
		).ServeHTTP(w, r)
		return
	}

	if err := c.paymentService.Create(r.Context(), &dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.Redirect(w, r, c.basePath)
}
