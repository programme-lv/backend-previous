package lang

import (
	"github.com/programme-lv/backend/internal/comm/i18nerror"
	"golang.org/x/text/language"
)

func newErrorLanguageNotFound() i18nerror.I18NError {
	return i18nerror.New("err_language_not_found", map[language.Tag]string{
		language.English: "language not found",
		language.Latvian: "valoda nav atrasta",
	}, nil)
}
