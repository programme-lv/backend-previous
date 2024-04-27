package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"strings"
)

func getGQLReqLang(ctx context.Context) string {
	opCtx := graphql.GetOperationContext(ctx)
	lang_row := opCtx.Headers.Get("Accept-Language")
	// split by comma
	langs := strings.Split(lang_row, ",")
	// split by semicolon
	lang_with_q := strings.Split(langs[0], ";")
	lang := lang_with_q[0]
	if lang == "lv-LV" {
		return "lv"
	}
	return lang
}
