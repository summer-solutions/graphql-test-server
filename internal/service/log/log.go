package log

import (
	"context"
	"summer-solutions/graphql-test-server/internal/gin"

	"github.com/apex/log"
)

const contextKey = "_log"

var logProvider Provider

func FromContext(ctx context.Context) log.Interface {
	g := gin.FromContext(ctx)
	l, has := g.Get(contextKey)
	if !has {
		if logProvider != nil {
			l = logProvider(ctx)
		} else {
			l = log.WithFields(&log.Fields{})
		}
		g.Set(contextKey, l)
	}
	return l.(*log.Entry)
}

type Provider func(ctx context.Context) log.Interface

func SetProvider(provider Provider) {
	logProvider = provider
}
