package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/programme-lv/backend/internal/i18nerror"
	"golang.org/x/text/language"
)

func newErrorInternalServer() i18nerror.I18NError {
	return i18nerror.New("internal_server_error", map[language.Tag]string{
		language.English: "internal server error",
		language.Latvian: "iekšējā servera kļūda",
	}, nil)
}

func newErrorUnauthorized() i18nerror.I18NError {
	return i18nerror.New("unauthorized", map[language.Tag]string{
		language.English: "unauthorized",
		language.Latvian: "neatļauts",
	}, nil)
}

// smartError filters out the error and returns an i18n error if it is one, otherwise returns a new internal server error.
func smartError(ctx context.Context, err error) i18nerror.I18NError {
	var i18nErr i18nerror.I18NError
	if errors.As(err, &i18nErr) {
		return i18nErr.WithLanguage(getGQLRequestLanguage(ctx))
	}
	return newErrorInternalServer().WithLanguage(getGQLRequestLanguage(ctx))
}

func getGQLRequestLanguage(ctx context.Context) string {
	opCtx := graphql.GetOperationContext(ctx)
	langRow := opCtx.Headers.Get("Accept-Language")
	return langRow
}
