package controllers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/presentation/templates/pages/error_pages"
	"github.com/iota-agency/iota-sdk/pkg/types"
)

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageCtx, err := composables.UsePageCtx(r, types.NewPageData("ErrorPages.NotFound.Title", ""))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		props := &error_pages.NotFoundPageProps{
			PageContext: pageCtx,
		}
		templ.Handler(error_pages.NotFoundContent(props), templ.WithStreaming()).ServeHTTP(w, r)
	}
}

func MethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
