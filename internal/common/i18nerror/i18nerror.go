package i18nerror

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
)

type I18NError interface {
	Error() string
	WithLanguage(lang string) I18NError
	Localize(languages ...string) string
}

type i18NErrorImpl struct {
	I18NMessageID  string
	Translations   *map[language.Tag]string
	ChosenLanguage *string
	TemplateData   *map[string]interface{}
}

func (e i18NErrorImpl) Error() string {
	if e.ChosenLanguage != nil {
		return e.Localize(*e.ChosenLanguage)
	} else {
		return e.Localize("en")
	}
}

func (e i18NErrorImpl) WithLanguage(lang string) I18NError {
	return &i18NErrorImpl{
		I18NMessageID:  e.I18NMessageID,
		Translations:   e.Translations,
		ChosenLanguage: &lang,
		TemplateData:   e.TemplateData,
	}
}

func (e i18NErrorImpl) Localize(langs ...string) string {
	bundle := i18n.NewBundle(language.English)
	for lang, translation := range *e.Translations {
		err := bundle.AddMessages(lang, &i18n.Message{
			ID:    e.I18NMessageID,
			Other: translation,
		})
		if err != nil {
			slog.Error(err.Error())
			return ""
		}
	}

	return i18n.NewLocalizer(bundle, langs...).MustLocalize(&i18n.LocalizeConfig{
		MessageID:    e.I18NMessageID,
		TemplateData: *e.TemplateData,
	})
}

func New(errorID string, translations map[language.Tag]string, templateData map[string]interface{}) I18NError {
	return i18NErrorImpl{
		I18NMessageID:  errorID,
		Translations:   &translations,
		ChosenLanguage: nil,
		TemplateData:   &templateData,
	}
}
