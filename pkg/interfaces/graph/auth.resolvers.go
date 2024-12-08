package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"

	model "github.com/iota-agency/iota-sdk/pkg/interfaces/graph/gqlmodels"
)

// Authenticate is the resolver for the authenticate field.
func (r *mutationResolver) Authenticate(ctx context.Context, email string, password string) (*model.Session, error) {
	//writer, ok := composables.UseWriter(ctx)
	//if !ok {
	//	return nil, fmt.Errorf("request params not found")
	//}
	//_, session, err := r.app.AuthService.Authenticate(ctx, email, password)
	//if err != nil {
	//	return nil, err
	//}
	//conf := configuration.Use()
	//cookie := &http.Cookie{
	//	Path:     conf.SidCookieKey,
	//	Value:    session.Token,
	//	Expires:  session.ExpiresAt,
	//	HttpOnly: false,
	//	SameSite: http.SameSiteDefaultMode,
	//	Secure:   false,
	//	Domain:   conf.Domain,
	//}
	//http.SetCookie(writer, cookie)
	//return session.ToGraph(), nil
	panic(fmt.Errorf("not implemented: Authenticate - authenticate"))
}

// GoogleAuthenticate is the resolver for the googleAuthenticate field.
func (r *mutationResolver) GoogleAuthenticate(ctx context.Context) (string, error) {
	//return r.app.AuthService.GoogleAuthenticate()
	panic(fmt.Errorf("not implemented: GoogleAuthenticate - googleAuthenticate"))
}

// DeleteSession is the resolver for the deleteSession field.
func (r *mutationResolver) DeleteSession(ctx context.Context, token string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteSession - deleteSession"))
}

// AuthenticationLog is the resolver for the authenticationLog field.
func (r *queryResolver) AuthenticationLog(ctx context.Context, id int64) (*model.AuthenticationLog, error) {
	panic(fmt.Errorf("not implemented: AuthenticationLog - authenticationLog"))
}

// AuthenticationLogs is the resolver for the authenticationLogs field.
func (r *queryResolver) AuthenticationLogs(ctx context.Context, offset int, limit int, sortBy []string) (*model.PaginatedAuthenticationLogs, error) {
	panic(fmt.Errorf("not implemented: AuthenticationLogs - authenticationLogs"))
}

// Session is the resolver for the session field.
func (r *queryResolver) Session(ctx context.Context, token string) (*model.Session, error) {
	panic(fmt.Errorf("not implemented: Session - session"))
}

// Sessions is the resolver for the sessions field.
func (r *queryResolver) Sessions(ctx context.Context, offset int, limit int, sortBy []string) (*model.PaginatedSessions, error) {
	panic(fmt.Errorf("not implemented: Sessions - sessions"))
}

// SessionDeleted is the resolver for the sessionDeleted field.
func (r *subscriptionResolver) SessionDeleted(ctx context.Context) (<-chan int64, error) {
	panic(fmt.Errorf("not implemented: SessionDeleted - sessionDeleted"))
}
