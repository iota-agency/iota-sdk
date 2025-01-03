package composables

import (
	"errors"
	"github.com/iota-uz/iota-sdk/pkg/constants"
	"github.com/iota-uz/iota-sdk/pkg/types"
	"net/http"
)

var (
	ErrNavItemsNotFound = errors.New("navigation items not found")
)

func UseNavItems(r *http.Request) ([]types.NavigationItem, error) {
	navItems := r.Context().Value(constants.NavItemsKey)
	if navItems == nil {
		return nil, ErrNavItemsNotFound
	}
	return navItems.([]types.NavigationItem), nil
}

func UseAllNavItems(r *http.Request) ([]types.NavigationItem, error) {
	navItems := r.Context().Value(constants.AllNavItemsKey)
	if navItems == nil {
		return nil, ErrNavItemsNotFound
	}
	return navItems.([]types.NavigationItem), nil
}
