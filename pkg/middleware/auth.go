package middleware

import (
	"context"
	services2 "github.com/iota-agency/iota-sdk/modules/core/services"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iota-agency/iota-sdk/pkg/configuration"
	"github.com/iota-agency/iota-sdk/pkg/constants"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/session"
)

func Authorize() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				conf := configuration.Use()
				token, err := r.Cookie(conf.SidCookieKey)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				ctx := r.Context()
				app, err := composables.UseApp(ctx)
				if err != nil {
					panic(err)
				}
				authService := app.Service(services2.AuthService{}).(*services2.AuthService)
				sess, err := authService.Authorize(ctx, token.Value)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				params, ok := composables.UseParams(ctx)
				if !ok {
					panic("params not found. Add RequestParams middleware up the chain")
				}
				params.Authenticated = true
				logger, err := composables.UseLogger(r.Context())
				if err == nil {
					logger.WithField("duration", time.Since(start)).Info("middleware.Authorize")
				}
				next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, constants.SessionKey, sess)))
			},
		)
	}
}

func ProvideUser() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				sess, ok := ctx.Value(constants.SessionKey).(*session.Session)
				if !ok {
					next.ServeHTTP(w, r)
					return
				}
				app, err := composables.UseApp(ctx)
				if err != nil {
					panic(err)
				}
				userService := app.Service(services2.UserService{}).(*services2.UserService)
				u, err := userService.GetByID(ctx, sess.UserID)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				ctx = context.WithValue(ctx, constants.UserKey, u)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func RedirectNotAuthenticated() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				params, ok := composables.UseParams(r.Context())
				if !ok {
					panic("params not found. Add RequestParams middleware up the chain")
				}
				if !params.Authenticated {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}

func RequireAuthorization() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				params, ok := composables.UseParams(r.Context())
				if !ok {
					panic("params not found. Add RequestParams middleware up the chain")
				}
				if !params.Authenticated {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}
