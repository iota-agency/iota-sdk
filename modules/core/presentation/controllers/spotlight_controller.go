package controllers

import (
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/iota-agency/iota-sdk/components/spotlight"
	"github.com/iota-agency/iota-sdk/modules/core/services"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/middleware"
	"github.com/iota-agency/iota-sdk/pkg/presentation/templates/icons"
	"github.com/iota-agency/iota-sdk/pkg/types"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"sort"
)

func flatNavItems(items []types.NavigationItem) []types.NavigationItem {
	var flatItems []types.NavigationItem
	for _, item := range items {
		flatItems = append(flatItems, item)
		if item.Children != nil {
			for _, child := range flatNavItems(item.Children) {
				flatItems = append(flatItems, types.NavigationItem{
					Name:     child.Name,
					Href:     child.Href,
					Icon:     item.Icon,
					Children: child.Children,
				})
			}
		}
	}
	return flatItems
}

type SpotlightController struct {
	app        application.Application
	tabService *services.TabService
	basePath   string
}

func NewSpotlightController(app application.Application) application.Controller {
	return &SpotlightController{
		app:        app,
		tabService: app.Service(services.TabService{}).(*services.TabService),
		basePath:   "/spotlight",
	}
}

func (c *SpotlightController) Key() string {
	return c.basePath
}

func (c *SpotlightController) Register(r *mux.Router) {
	router := r.PathPrefix(c.basePath).Subrouter()
	router.Use(
		middleware.Authorize(),
		middleware.ProvideUser(),
		middleware.RedirectNotAuthenticated(),
		middleware.WithLocalizer(c.app.Bundle()),
	)
	router.HandleFunc("/search", c.Get).Methods(http.MethodGet)
}

func (c *SpotlightController) spotlightItems(localizer *i18n.Localizer) []*spotlight.SpotlightItem {
	navItems := flatNavItems(c.app.NavItems(localizer))
	t := func(id string) string {
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: id,
		})
	}
	items := []*spotlight.SpotlightItem{
		{
			Icon:  icons.PlusCircle(icons.Props{Size: "24"}),
			Title: t("Spotlight.Actions.AddUser"),
			Link:  "/users/new",
		},
		{
			Icon:  icons.PlusCircle(icons.Props{Size: "24"}),
			Title: t("Spotlight.Actions.AddExpense"),
			Link:  "/finance/expenses/new",
		},
		{
			Icon:  icons.PlusCircle(icons.Props{Size: "24"}),
			Title: t("Spotlight.Actions.AddPayment"),
			Link:  "/finance/payments/new",
		},
		{
			Icon:  icons.PlusCircle(icons.Props{Size: "24"}),
			Title: t("Spotlight.Actions.AddAccount"),
			Link:  "/finance/accounts/new",
		},
		{
			Icon:  icons.SignOut(icons.Props{Size: "24"}),
			Title: t("NavigationLinks.Navbar.Logout"),
			Link:  "/logout",
		},
	}
	for _, item := range navItems {
		items = append(items, &spotlight.SpotlightItem{
			Icon:  item.Icon,
			Title: item.Name,
			Link:  item.Href,
		})
	}
	return items
}

func (c *SpotlightController) Get(w http.ResponseWriter, r *http.Request) {
	localizer, ok := composables.UseLocalizer(r.Context())
	if !ok {
		http.Error(w, composables.ErrNoLocalizer.Error(), http.StatusInternalServerError)
		return
	}
	q := r.URL.Query().Get("q")
	if q == "" {
		templ.Handler(spotlight.SpotlightItems([]*spotlight.SpotlightItem{})).ServeHTTP(w, r)
		return
	}

	items := c.spotlightItems(localizer)
	var filteredItems []*spotlight.SpotlightItem
	words := make([]string, len(items))
	for i, item := range items {
		words[i] = item.Title
	}
	ranks := fuzzy.RankFindNormalizedFold(q, words)
	sort.Sort(ranks)
	for _, rank := range ranks {
		filteredItems = append(filteredItems, items[rank.OriginalIndex])
	}
	templ.Handler(spotlight.SpotlightItems(filteredItems)).ServeHTTP(w, r)
}
