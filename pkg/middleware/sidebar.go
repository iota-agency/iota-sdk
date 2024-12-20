package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/constants"
	"github.com/iota-agency/iota-sdk/pkg/domain/aggregates/user"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/tab"
	"github.com/iota-agency/iota-sdk/pkg/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

func filterItems(items []types.NavigationItem, user *user.User) []types.NavigationItem {
	filteredItems := make([]types.NavigationItem, 0, len(items))
	for _, item := range items {
		if item.HasPermission(user) {
			filteredItems = append(filteredItems, types.NavigationItem{
				Name:        item.Name,
				Href:        item.Href,
				Children:    filterItems(item.Children, user),
				Icon:        item.Icon,
				Permissions: item.Permissions,
			})
		}
	}
	return filteredItems
}

func hrefExists(href string, tabs []*tab.Tab) bool {
	for _, tab := range tabs {
		if tab.Href == href {
			return true
		}
	}
	return false
}

func getEnabledNavItems(items []types.NavigationItem, tabs []*tab.Tab) []types.NavigationItem {
	var out []types.NavigationItem
	for _, item := range items {
		if len(item.Children) > 0 {
			children := getEnabledNavItems(item.Children, tabs)
			childrenLen := len(children)
			if childrenLen == 0 {
				continue
			}
			if childrenLen == 1 {
				out = append(out, children[0])
			} else {
				item.Children = children
				out = append(out, item)
			}
		} else if item.Href == "" || item.Href != "" && hrefExists(item.Href, tabs) {
			out = append(out, item)
		}
	}

	return out
}

func getNavItems(
	app application.Application,
	localizer *i18n.Localizer,
	user *user.User,
	tabs []*tab.Tab,
) ([]types.NavigationItem, []types.NavigationItem) {
	filtered := filterItems(app.NavigationItems(localizer), user)
	return getEnabledNavItems(filtered, tabs), filtered
}

func NavItems(app application.Application) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				localizer, ok := composables.UseLocalizer(r.Context())
				if !ok {
					panic("localizer not found")
				}
				u, err := composables.UseUser(r.Context())
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				tabs, err := composables.UseTabs(r.Context())
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				items, allItems := getNavItems(app, localizer, u, tabs)
				ctx := context.WithValue(r.Context(), constants.NavItemsKey, items)
				ctx = context.WithValue(ctx, constants.AllNavItemsKey, allItems)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
