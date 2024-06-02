package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

func getGQLRequestLanguage(ctx context.Context) string {
	opCtx := graphql.GetOperationContext(ctx)
	lang_row := opCtx.Headers.Get("Accept-Language")
	return lang_row
}
