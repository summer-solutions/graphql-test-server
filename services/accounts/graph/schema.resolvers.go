package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"summer-solutions/graphql-test-server/services/accounts/graph/generated"
	"summer-solutions/graphql-test-server/services/accounts/graph/model"
	"time"
)

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	fmt.Printf("%v\n", ctx.Value("www").(time.Time).String())
	return &model.User{
		ID:       "1234",
		Username: "Me",
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
