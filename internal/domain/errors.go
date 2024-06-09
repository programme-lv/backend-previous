package domain

import (
	"github.com/programme-lv/backend/internal/i18nerror"
	"golang.org/x/text/language"
)

func NewErrorInternalServer() i18nerror.I18NError {
	return i18nerror.New("err_internal_server", map[language.Tag]string{
		language.English: "internal server error",
		language.Latvian: "iekšējā servera kļūda",
	}, nil)
}
