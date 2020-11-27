package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"summer-solutions/graphql-test-server/internal/gin"
	"summer-solutions/graphql-test-server/pkg/service"
	"summer-solutions/graphql-test-server/services/accounts/graph/generated"
	"summer-solutions/graphql-test-server/services/accounts/graph/model"
)

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	fmt.Println(gin.FromContext(ctx).Get(service.OrmContextService))
	return &model.User{
		ID:       "1234",
		Username: "Me",
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
