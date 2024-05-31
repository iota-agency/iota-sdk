package positions

import (
	"github.com/graphql-go/graphql"
	"github.com/iota-agency/iota-erp/models"
	"github.com/iota-agency/iota-erp/sdk/graphql/old/resolvers"
	"gorm.io/gorm"
)

func Queries(db *gorm.DB) []*graphql.Field {
	return resolvers.DefaultQueries(db, &models.Position{}, "position", "positions")
}

func Mutations(db *gorm.DB) []*graphql.Field {
	return resolvers.DefaultMutations(db, &models.Position{}, "position")
}
