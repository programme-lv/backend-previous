package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

func getGQLReqLang(ctx context.Context) string {
	opCtx := graphql.GetOperationContext(ctx)
	lang := opCtx.Headers.Get("Accept-Language")
	return lang
}
